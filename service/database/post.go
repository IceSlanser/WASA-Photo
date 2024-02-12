package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) GetPosts(myUID uint64, userID uint64) ([]Post, error) {
	// Store posts
	rows, err := db.c.Query("SELECT * FROM posts WHERE ProfileID = ? AND ProfileID NOT IN (SELECT BannerUID FROM bans WHERE BannedUID = ?)", userID, myUID)
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

func (db *appdbimpl) PostPost(UID uint64, photo []byte, description string) (uint64, error) {
	// Post a new post
	res, err := db.c.Exec("INSERT INTO posts (ProfileID, File, Description) VALUES (?, ?, ?)", UID, photo, description)
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

func (db *appdbimpl) GetLikes(myUID uint64, postID uint64) ([]uint64, uint64, error) {
	// Store likes
	query := `SELECT OwnerID 
				FROM likes
				WHERE PostID = ? AND OwnerID NOT IN (SELECT BannerUID FROM bans WHERE BannedUID = ?)`
	rows, err := db.c.Query(query, postID, myUID)
	var likeOwners []uint64
	for rows.Next() {
		var ownerID uint64
		err = rows.Scan(&ownerID)
		if err != nil {
			return nil, 0, err
		}
		likeOwners = append(likeOwners, ownerID)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	// Get LikeCount
	var likeCount uint64
	err = db.c.QueryRow("SELECT LikeCount FROM posts WHERE ID = ?", postID).Scan(&likeCount)
	if err != nil {
		return nil, 0, err
	}

	return likeOwners, likeCount, nil
}

func (db *appdbimpl) PutLike(UID uint64, postID uint64) (uint64, bool, error) {
	valid, err := db.IsValidPost(postID)
	if err != nil {
		return 0, false, err
	}
	if !valid {
		return 0, false, err
	}

	// Try to insert a like
	var like Like
	res, err := db.c.Exec("INSERT INTO likes(PostID, OwnerID) VALUES (?, ?)",
		postID, UID)
	if err != nil {
		err = db.c.QueryRow("SELECT * FROM likes WHERE PostID = ? AND OwnerID = ?",
			postID, UID).Scan(&like.ID, &like.PostID, &like.OwnerID)
		// There is already an existent like
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return 0, false, err
			}
		}
		return like.ID, true, nil
	}

	// Get the new like ID
	ID, err := res.LastInsertId()
	if err != nil {
		return like.ID, false, err
	}
	like.ID = uint64(ID)

	// Update LikeCount
	_, err = db.c.Exec("UPDATE posts SET LikeCount = LikeCount + 1 WHERE ID = ?", postID)
	if err != nil {
		return like.ID, true, err
	}

	return like.ID, false, nil
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

func (db *appdbimpl) GetComments(myUID uint64, postID uint64) ([]Comment, uint64, error) {
	// Store comments
	query := `SELECT * 
				FROM comments 
				WHERE PostID = ? AND OwnerID NOT IN (SELECT BannerUID FROM bans WHERE BannedUID = ?)`
	rows, err := db.c.Query(query, postID, myUID)
	var comments []Comment
	for rows.Next() {
		var comment Comment
		err = rows.Scan(&comment.ID, &comment.PostID, &comment.OwnerID, &comment.Text, &comment.DateTime)
		if err != nil {
			return nil, 0, err
		}
		comments = append(comments, comment)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	// Get CommentCount
	var CommentCount uint64
	err = db.c.QueryRow("SELECT CommentCount FROM posts WHERE ID = ?", postID).Scan(&CommentCount)
	if err != nil {
		return nil, 0, err
	}

	return comments, CommentCount, nil
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

	// Check if the user are allowed to delete the comment
	if ownerID == UID || ownerID == userID {
		_, err = db.c.Exec("DELETE FROM comments WHERE ID = ?", commentID)
		if err != nil {
			return true, err
		}
		return true, nil
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
