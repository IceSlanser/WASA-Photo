package database

import (
	"database/sql"
	"time"
)

func (db *appdbimpl) LoginUser(name string) (uint64, bool, error) {
	var ID uint64

	// Try to insert the user into the database
	res, err := db.c.Exec("INSERT INTO profiles(Username) VALUES (?)", name)
	if err != nil {
		err := db.c.QueryRow("SELECT ID FROM profiles WHERE Username = ?", name).Scan(&ID)
		if err != nil {
			// There is already an existent user with the input username
			if err == sql.ErrNoRows {
				return 0, true, err
			}
		}
		return ID, true, nil
	}

	// A new user has been created
	UID, err := res.LastInsertId()
	if err != nil {
		return 0, false, err
	}
	return uint64(UID), false, nil
}

func (db *appdbimpl) SetUsername(UID uint64, nname string) error {
	// Update your username with a new one
	_, err := db.c.Exec("UPDATE profiles SET Username = ? WHERE ID = ?", nname, UID)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) GetProfile(myUID uint64, userID uint64) (User, error) {
	var user User
	valid, err := db.IsValidProfile(userID)
	if err != nil {
		return User{}, err
	}
	if !valid {
		return User{}, err
	}

	err = db.c.QueryRow("SELECT * FROM profiles WHERE ID = ? AND ID NOT IN (SELECT BannerUID FROM bans WHERE BannedUID = ?)", userID, myUID).Scan(&user.ID, &user.Username,
		&user.FollowingCount, &user.FollowersCount, &user.PostCount)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, err
		}
	}

	return user, nil
}

func (db *appdbimpl) GetStream(UID uint64, startTime time.Time, endTime time.Time) ([]Post, error) {
	query := `SELECT posts.*
				FROM posts
				LEFT JOIN follows ON FollowedUID = ProfileID
				WHERE ProfileID NOT IN (SELECT BannerUID FROM bans WHERE BannedUID = ?) AND follows.FollowerUID = ? AND DateTime BETWEEN ? AND ?
			  	ORDER BY DateTime DESC`
	rows, err := db.c.Query(query, UID, UID, startTime, endTime)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Store posts
	var posts []Post
	for rows.Next() {
		var post Post
		err = rows.Scan(&post.ID, &post.ProfileID, &post.File, &post.Description, &post.LikeCount, &post.CommentCount, &post.DateTime)
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

func (db *appdbimpl) GetFollows(myUID uint64, userID uint64) ([]uint64, []uint64, error) {
	// Store followings
	query := `SELECT FollowedUID 
				FROM follows 
				WHERE FollowerUID = ? AND FollowedUID NOT IN (SELECT BannerUID FROM bans WHERE BannedUID = ?)`
	followingRows, err := db.c.Query(query, userID, myUID)
	var followings []uint64
	for followingRows.Next() {
		var followingUID uint64
		err = followingRows.Scan(&followingUID)
		if err != nil {
			return nil, nil, err
		}
		followings = append(followings, followingUID)
	}
	if err := followingRows.Err(); err != nil {
		return nil, nil, err
	}

	// Store followers
	query2 := `SELECT FollowerUID 
				FROM follows 
				WHERE FollowedUID = ? AND FollowedUID NOT IN (SELECT BannerUID FROM bans WHERE BannedUID = ?)`
	followerRows, err := db.c.Query(query2, userID, myUID)
	var followers []uint64
	for followerRows.Next() {
		var followerUID uint64
		err = followerRows.Scan(&followerUID)
		if err != nil {
			return nil, nil, err
		}
		followers = append(followers, followerUID)
	}
	if err := followerRows.Err(); err != nil {
		return nil, nil, err
	}

	return followings, followers, nil
}

func (db *appdbimpl) PutFollow(UID uint64, followedUID uint64) (uint64, bool, error) {
	valid, err := db.IsValidProfile(followedUID)
	if err != nil {
		return 0, false, err
	}
	if !valid {
		return 0, false, err
	}

	var followID uint64
	// Try to insert the follow into the database
	res, err := db.c.Exec("INSERT INTO follows(FollowerUID, FollowedUID) VALUES (?, ?)", UID, followedUID)
	if err != nil {
		err = db.c.QueryRow("SELECT ID FROM follows WHERE FollowerUID = ? AND FollowedUID = ?", UID, followedUID).Scan(&followID)
		if err != nil {
			// There is already an existent follow
			if err == sql.ErrNoRows {
				return 0, false, err
			}
		}
		return followID, true, nil
	}

	// A new follow has been added
	ID, err := res.LastInsertId()
	if err != nil {
		return 0, false, err
	}

	// Update FollowingCount for the follower
	_, err = db.c.Exec("UPDATE profiles SET FollowingCount = FollowingCount + 1 WHERE ID = ?", UID)
	if err != nil {
		return 0, false, err
	}

	// Update FollowerCount for the followed user
	_, err = db.c.Exec("UPDATE profiles SET FollowerCount = FollowerCount + 1 WHERE ID = ?", followedUID)
	if err != nil {
		return 0, false, err
	}

	return uint64(ID), true, nil
}

func (db *appdbimpl) DeleteFollow(UID uint64, followedUID uint64) (bool, error) {
	// Check if there is an existent follow
	var followID uint64
	err := db.c.QueryRow("SELECT ID FROM follows WHERE FollowerUID = ? AND FollowedUID = ?", UID, followedUID).Scan(&followID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, err
		}
	}

	// Unfollow
	_, err = db.c.Exec("DELETE FROM follows WHERE ID = ?", followID)
	if err != nil {
		return true, err
	}

	// Update FollowingCount for the follower
	_, err = db.c.Exec("UPDATE profiles SET FollowingCount = FollowingCount - 1 WHERE ID = ?", UID)
	if err != nil {
		return true, err
	}

	// Update FollowerCount for the followed user
	_, err = db.c.Exec("UPDATE profiles SET FollowerCount = FollowerCount - 1 WHERE ID = ?", followedUID)
	if err != nil {
		return true, err
	}

	return true, nil
}

func (db *appdbimpl) PutBan(UID uint64, bannedUID uint64) (uint64, bool, error) {
	valid, err := db.IsValidProfile(bannedUID)
	if err != nil {
		return 0, false, err
	}
	if !valid {
		return 0, false, err
	}

	var banID uint64
	// Try to insert the follow into the database
	res, err := db.c.Exec("INSERT INTO bans(BannerUID, BannedUID) VALUES (?, ?)", UID, bannedUID)
	if err != nil {
		err = db.c.QueryRow("SELECT ID FROM bans WHERE BannerUID = ? AND BannedUID = ?", UID, bannedUID).Scan(&banID)
		if err != nil {
			// There is already an existent follow
			if err == sql.ErrNoRows {
				return 0, false, err
			}
		}
		return banID, true, nil
	}

	// A new follow has been added
	ID, err := res.LastInsertId()
	if err != nil {
		return 0, false, err
	}
	return uint64(ID), false, nil
}

func (db *appdbimpl) DeleteBan(UID uint64, bannedUID uint64) (bool, error) {
	// Check if there is an existent follow
	var banID uint64
	err := db.c.QueryRow("SELECT ID FROM bans WHERE BannerUID = ? AND BannedUID = ?", UID, bannedUID).Scan(&banID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, err
		}
	}

	// Unban
	_, err = db.c.Exec("DELETE FROM bans WHERE ID = ?", banID)
	if err != nil {
		return true, err
	}
	return true, nil
}

func (db *appdbimpl) IsAvailable(newname string) (uint64, bool) {
	var UID uint64

	// Return true if the username is not taken, false otherwise
	err := db.c.QueryRow("SELECT ID FROM profiles WHERE Username = ?", newname).Scan(&UID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, true
		}
	}
	return UID, false
}

func (db *appdbimpl) IsValidProfile(ID uint64) (bool, error) {
	var foo uint64
	err := db.c.QueryRow("SELECT ID FROM profiles WHERE ID = ?", ID).Scan(&foo)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, err
		}
	}
	return true, nil
}
