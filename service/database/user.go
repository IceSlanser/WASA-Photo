package database

import (
	"database/sql"
	"time"
)

func (db *appdbimpl) LoginUser(name string) (uint64, bool, error) {
	var u User
	u.Username = name

	// Try to insert the user into the database
	res, err := db.c.Exec("INSERT INTO profiles(Username) VALUES (?)", u.Username)
	if err != nil {
		err := db.c.QueryRow("SELECT ID, Username FROM profiles WHERE Username = ?", u.Username).Scan(&u.ID)
		if err != nil {
			// There is already an existent user with the input username
			if err == sql.ErrNoRows {
				return u.ID, false, ErrorUserDoesNotExist
			}
		}
		return u.ID, false, nil
	}

	// A new user has been created
	ID, err := res.LastInsertId()
	if err != nil {
		return u.ID, false, err
	}
	u.ID = uint64(ID)
	return u.ID, true, nil
}

func (db *appdbimpl) SetUsername(UID uint64, nname string) error {
	// Update your username with a new one
	_, err := db.c.Exec("UPDATE profiles SET Username = ? WHERE ID = ?", nname, UID)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) IsAvailable(newname string) (bool, error) {
	var username string

	// Return true if the username is not taken, false otherwise
	err := db.c.QueryRow("SELECT Username FROM profiles WHERE Username = ?", newname).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, err
		}
	}
	return false, nil
}

func (db *appdbimpl) GetProfile(ID uint64) (User, error) {
	var user User

	/*
		valid, err := db.IsValid(ID)
		if err != nil {
			return User{}, err
		}
		if !valid {
			return User{}, err
		}
	*/

	err := db.c.QueryRow("SELECT * FROM profiles WHERE ID = ?", ID).Scan(&user.ID, &user.Username,
		&user.Following, &user.Followers, &user.Post)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, err
		}
	}

	return user, nil
}

func (db *appdbimpl) IsValid(ID uint64) (bool, error) {
	var foo uint64
	err := db.c.QueryRow("SELECT ID FROM profiles WHERE ID = ?", ID).Scan(&foo)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, err
		}
	}
	return true, nil
}

func (db *appdbimpl) GetFollowing(ID uint64) ([]uint64, error) {
	// Select all followings
	rows, err := db.c.Query("SELECT FollowedUID FROM follows WHERE FollowerUID = ?", ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Store following
	var following []uint64
	for rows.Next() {
		var followingUID uint64
		if err := rows.Scan(&followingUID); err != nil {
			return nil, err
		}
		following = append(following, followingUID)
	}
	return following, nil
}

func (db *appdbimpl) GetPosts(UID uint64, startTime time.Time, endTime time.Time) ([]Post, error) {
	query := `SELECT posts.*, EXISTS(SELECT * FROM bans WHERE BannerUID = ProfileID and BannedUID = 1) as banned
				FROM posts
				LEFT JOIN follows ON FollowedUID = ProfileID
				WHERE follows.FollowerUID = ? AND DateTime BETWEEN ? AND ?
			  	ORDER BY DateTime DESC`
	//////////////////// Vedere se esiste un modo migliore per vedere i ban
	rows, err := db.c.Query(query, UID, UID, startTime, endTime)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Store posts
	var posts []Post
	for rows.Next() {
		var post Post
		/// IN caso controllare la colonna banned
		err = rows.Scan(&post.ID, &post.ProfileID, &post.Description, &post.LikeCount, &post.CommentCount, &post.DateTime)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (db *appdbimpl) PutFollow(followedUID uint64, UID uint64) (uint64, bool, error) {
	var followerFoo uint64
	var followedFoo uint64

	// Try to insert the follow into the database
	res, err := db.c.Exec("INSERT INTO follows(followerUID, followedUID) VALUES (?, ?)", UID, followedUID)
	if err != nil {
		err = db.c.QueryRow("SELECT * FROM follows WHERE FollowerUID = ? AND FollowedUID = ?", followerFoo, followedFoo).Scan()
		if err != nil {
			// There is already an existent follow
			if err == sql.ErrNoRows {
				return 0, false, err
			}
		}
		return 0, true, nil
	}

	// A new follow has been added
	ID, err := res.LastInsertId()
	if err != nil {
		return 0, false, err
	}
	return uint64(ID), false, nil
}

func (db *appdbimpl) DeleteFollow(UID uint64, followedUID uint64) (bool, error) {
	// Check if there is an existent follow
	err := db.c.QueryRow("SELECT ID FROM follows WHERE FollowerUID = ? AND FollowedUID = ?", UID, followedUID).Scan()
	if err != nil {
		if err == sql.ErrNoRows {
			return false, err
		}
	}

	// Check if the user are allowed to delete that like
	var fooUID uint64
	if fooUID != UID { //////////////Passare inpath username e non ID
		return false, nil
	}
	_, err = db.c.Exec("DELETE FROM likes WHERE ID = ?", followedUID)
	if err != nil {
		return true, err
	}
	return true, nil
}
