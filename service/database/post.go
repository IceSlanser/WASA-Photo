package database

import (
	"database/sql"
	"errors"
	"time"
)

func (db *appdbimpl) GetUserPosts(myUID uint64, userID uint64) ([]Post, error) {
	// Store posts
	rows, err := db.c.Query("SELECT * FROM posts WHERE ProfileID = ? AND ProfileID NOT IN (SELECT BannerUID FROM bans WHERE BannedUID = ?)", userID, myUID)
	if err != nil {
		return nil, err
	}
	var posts []Post
	for rows.Next() {
		var post Post
		err = rows.Scan(&post.ID, &post.ProfileID, &post.File, &post.Description, &post.LikeCount, &post.CommentCount, &post.DateTime)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (db *appdbimpl) GetPostInfo(myUID uint64, postID uint64) (Post, error) {
	// Store posts
	var post Post
	err := db.c.QueryRow("SELECT * FROM posts WHERE ID = ? AND ProfileID NOT IN (SELECT BannerUID FROM bans WHERE BannedUID = ?)", postID, myUID).Scan(&post.ID, &post.ProfileID, &post.File, &post.Description, &post.LikeCount, &post.CommentCount, &post.DateTime)
	if err != nil {
		return Post{}, err
	}

	return post, nil
}

func (db *appdbimpl) PostPost(UID uint64, photo []byte, description string) (uint64, error) {
	// Post a new post
	now := time.Now()
	dateTime := now.Format("2002-11-20 20:20:20 GMT+2")
	res, err := db.c.Exec("INSERT INTO posts (ProfileID, File, Description, DateTime) VALUES (?, ?, ?, ?)", UID, photo, description, dateTime)
	if err != nil {
		return 0, err
	}

	// Get a new post ID
	ID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Update PostCount
	_, err = db.c.Exec("UPDATE profiles SET PostCount = PostCount + 1 WHERE ID = ?", UID)
	if err != nil {
		return 0, err
	}

	return uint64(ID), nil
}

func (db *appdbimpl) DeletePost(UID uint64, postID uint64) (bool, error) {
	// Check if there is an existent post
	var ID uint64
	err := db.c.QueryRow("SELECT ID FROM posts WHERE ID = ? AND ProfileID = ?", postID, UID).Scan(&postID, &ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, err
		}
	}

	// Delete post
	_, err = db.c.Exec("DELETE FROM posts WHERE ID = ? AND ProfileID = ?", postID, UID)
	if err != nil {
		return true, err
	}

	// Update PostCount
	_, err = db.c.Exec("UPDATE profiles SET PostCount = PostCount - 1 WHERE ID = ?", UID)
	if err != nil {
		return true, err
	}

	return true, nil
}

func (db *appdbimpl) GetPhoto(UID uint64, postID uint64) ([]byte, error) {
	// Store photo
	var file []byte
	query := `SELECT File 
				FROM posts
				WHERE ID = ? AND ProfileID NOT IN (SELECT BannerUID FROM bans WHERE BannedUID = ?)`
	err := db.c.QueryRow(query, postID, UID).Scan(&file)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	return file, nil
}

func (db *appdbimpl) GetLikes(myUID uint64, postID uint64) ([]string, error) {
	// Store likes
	query := `SELECT OwnerID 
				FROM likes
				WHERE PostID = ? AND OwnerID NOT IN (SELECT BannerUID FROM bans WHERE BannedUID = ?)`
	rows, err := db.c.Query(query, postID, myUID)
	if err != nil {
		return nil, err
	}
	var likeOwners []string
	for rows.Next() {
		var ownerID uint64
		var ownerUsername string
		err = rows.Scan(&ownerID)
		if err != nil {
			return nil, err
		}
		err = db.c.QueryRow("SELECT Username FROM profiles WHERE ID = ?", ownerID).Scan(&ownerUsername)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, err
			}
		}
		likeOwners = append(likeOwners, ownerUsername)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return likeOwners, nil
}

func (db *appdbimpl) PutLike(UID uint64, postID uint64) (bool, error) {
	valid, err := db.IsValidPost(postID)
	if err != nil {
		return false, err
	}
	if !valid {
		return false, err
	}

	// Try to insert a like
	var likeID uint64
	_, err = db.c.Exec("INSERT INTO likes(PostID, OwnerID) VALUES (?, ?)",
		postID, UID)
	if err != nil {
		err = db.c.QueryRow("SELECT ID FROM likes WHERE PostID = ? AND OwnerID = ?",
			postID, UID).Scan(&likeID)
		// There is already an existent like
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return false, err
			}
		}
		return true, nil
	}

	// Update LikeCount
	_, err = db.c.Exec("UPDATE posts SET LikeCount = LikeCount + 1 WHERE ID = ?", postID)
	if err != nil {
		return true, err
	}

	return false, nil
}

func (db *appdbimpl) DeleteLike(UID uint64, postID uint64) (bool, error) {
	// Check if there is an existent like
	var likeID uint64
	err := db.c.QueryRow("SELECT ID FROM likes WHERE PostID = ? AND OwnerID = ?", postID, UID).Scan(&likeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, err
		}
	}

	_, err = db.c.Exec("DELETE FROM likes WHERE ID = ?", likeID)
	if err != nil {
		return true, err
	}

	// Update LikeCount
	_, err = db.c.Exec("UPDATE posts SET LikeCount = LikeCount - 1 WHERE ID = ?", postID)
	if err != nil {
		return true, err
	}

	return true, nil
}

func (db *appdbimpl) GetComments(myUID uint64, postID uint64) ([]Comment, error) {
	// Store comments
	query := `SELECT * 
				FROM comments 
				WHERE PostID = ? AND OwnerID NOT IN (SELECT BannerUID FROM bans WHERE BannedUID = ?)`
	rows, err := db.c.Query(query, postID, myUID)
	if err != nil {
		return nil, err
	}
	var comments []Comment
	for rows.Next() {
		var comment Comment
		err = rows.Scan(&comment.ID, &comment.PostID, &comment.OwnerID, &comment.Text, &comment.DateTime)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (db *appdbimpl) PostComment(UID uint64, postID uint64, text string) (uint64, error) {
	valid, err := db.IsValidPost(postID)
	if err != nil {
		return 0, err
	}
	if !valid {
		return 0, err
	}

	// Post new comment
	res, err := db.c.Exec("INSERT INTO comments(PostID, OwnerID, Text) VALUES (?, ?, ?)", postID, UID, text)
	if err != nil {
		return 0, err
	}

	// Get new comment ID
	ID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Update CommentCount
	_, err = db.c.Exec("UPDATE posts SET CommentCount = CommentCount + 1 WHERE ID = ?", postID)
	if err != nil {
		return uint64(ID), err
	}

	return uint64(ID), nil
}

func (db *appdbimpl) DeleteComment(UID uint64, postID uint64, commentID uint64) (bool, error) {
	// Check if there is an existent comment
	var ownerID uint64
	var userID uint64
	err := db.c.QueryRow("SELECT OwnerID FROM comments WHERE ID = ?", commentID).Scan(&ownerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, err
		}
	}

	// Find the post owner
	err = db.c.QueryRow("SELECT ProfileID FROM posts WHERE ID = ?", postID).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, err
		}
	}

	// Check if the user is allowed to delete the comment
	if ownerID == UID || ownerID == userID {
		_, err = db.c.Exec("DELETE FROM comments WHERE ID = ?", commentID)
		if err != nil {
			return true, err
		}
	}

	// Update CommentCount
	_, err = db.c.Exec("UPDATE posts SET CommentCount = CommentCount - 1 WHERE ID = ?", postID)
	if err != nil {
		return true, err
	}

	return true, nil
}

func (db *appdbimpl) IsValidPost(ID uint64) (bool, error) {
	var foo uint64
	err := db.c.QueryRow("SELECT ID FROM posts WHERE ID = ?", ID).Scan(&foo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, err
		}
	}
	return true, nil
}
