/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	LoginUser(string) (uint64, bool, error)
	SetUsername(uint64, string) error
	IsAvailable(string) bool

	GetProfile(uint64, uint64) (User, error)
	GetPosts(uint64, uint64) ([]Post, error)
	GetStream(uint64, time.Time, time.Time) ([]Post, error)
	GetFollows(uint64, uint64) ([]uint64, []uint64, error)
	GetComments(uint64, uint64) ([]Comment, error)
	GetLikes(uint64, uint64) ([]uint64, error)

	PostComment(uint64, uint64, string) (uint64, error)
	PostPost(uint64, []byte, string) (uint64, error)

	PutLike(uint64, uint64) (uint64, bool, error)
	PutFollow(uint64, uint64) (uint64, bool, error)
	PutBan(uint64, uint64) (uint64, bool, error)

	DeletePost(uint64, uint64) (bool, error)
	DeleteLike(uint64, uint64) (bool, error)
	DeleteComment(uint64, uint64) (bool, error)
	DeleteFollow(uint64, uint64) (bool, error)
	DeleteBan(uint64, uint64) (bool, error)
	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='example_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {

		// Database's tables
		profiles := `CREATE TABLE IF NOT EXISTS profiles (
			ID INTEGER PRIMARY KEY,
			Username TEXT UNIQUE,
			FollowingCount INTEGER DEFAULT 0,
			FollowerCount INTEGER DEFAULT 0,
			PostCount INTEGER DEFAULT 0
		)`

		posts := `CREATE TABLE IF NOT EXISTS posts (
			ID INTEGER PRIMARY KEY,
			ProfileID INTEGER NOT NULL,
			File []byte NOT NULL,
			Description TEXT NOT NULL,
			LikeCount INTEGER DEFAULT 0,
			CommentCount INTEGER DEFAULT 0,
			DateTime DATETIME DEFAULT CURRENT_TIMESTAMP,

			FOREIGN KEY (ProfileID) REFERENCES profiles(ID)
		)`

		follows := `CREATE TABLE IF NOT EXISTS follows (
			ID INTEGER PRIMARY KEY,
			FollowerUID INTEGER,
			FollowedUID INTEGER,
			UNIQUE (FollowerUID, FollowedUID),

			FOREIGN KEY (FollowerUID) REFERENCES profiles(ID) ON DELETE CASCADE,
			FOREIGN KEY (FollowedUID) REFERENCES profiles(ID) ON DELETE CASCADE
		)`

		bans := `CREATE TABLE IF NOT EXISTS bans (
			ID INTEGER PRIMARY KEY,
			BannerUID INTEGER,
			BannedUID INTEGER,
			UNIQUE (BannerUID, BannedUID),

			FOREIGN KEY (BannerUID) REFERENCES profiles(ID) ON DELETE CASCADE,
			FOREIGN KEY (BannedUID) REFERENCES profiles(ID) ON DELETE CASCADE
		)`

		likes := `CREATE TABLE IF NOT EXISTS likes (
			ID INTEGER PRIMARY KEY,
			PostID INTEGER,
			OwnerID INTEGER,
			UNIQUE (PostID, OwnerID),

			FOREIGN KEY (PostID) REFERENCES posts(ID) ON DELETE CASCADE,
			FOREIGN KEY (OwnerID) REFERENCES profiles(ID) ON DELETE CASCADE
		)`

		comments := `CREATE TABLE IF NOT EXISTS comments (
			ID INTEGER PRIMARY KEY,
			PostID INTEGER NOT NULL,
			OwnerID INTEGER NOT NULL,
			Text TEXT NOT NULL,
			DateTime DATETIME DEFAULT CURRENT_TIMESTAMP,

			FOREIGN KEY (PostID) REFERENCES posts(ID) ON DELETE CASCADE,
			FOREIGN KEY (OwnerID) REFERENCES profiles(ID) ON DELETE CASCADE
		)`

		// Tables error check
		_, err = db.Exec(profiles)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
		_, err = db.Exec(posts)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
		_, err = db.Exec(follows)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
		_, err = db.Exec(bans)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
		_, err = db.Exec(likes)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
		_, err = db.Exec(comments)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
