package database

import (
	"time"
)

type User struct {
	ID             uint64 `json:"ID"`
	Username       string `json:"Username"`
	FollowingCount uint64 `json:"Following_count"`
	FollowersCount uint64 `json:"Followers_count"`
	PostCount      uint64 `json:"Post_count"`
}

type Post struct {
	ID           uint64    `json:"ID"`
	ProfileID    string    `json:"profile_ID"`
	Description  string    `json:"description"`
	LikeCount    uint64    `json:"like_count"`
	CommentCount uint64    `json:"comment_count"`
	DateTime     time.Time `json:"date_time"`
}

type Like struct {
	ID      uint64 `json:"ID"`
	PostID  uint64 `json:"post_ID"`
	OwnerID uint64 `json:"owner_ID"`
}

type Comment struct {
	ID       uint64    `json:"ID"`
	PostID   uint64    `json:"post_ID"`
	OwnerID  uint64    `json:"owner_ID"`
	Text     string    `json:"text"`
	DateTime time.Time `json:"date_time"`
}
