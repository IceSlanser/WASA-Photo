package database

import (
	"time"
)

type User struct {
	ID             uint64 `json:"ID"`
	Username       string `json:"username"`
	FollowingCount uint64 `json:"following_count"`
	FollowerCount  uint64 `json:"follower_count"`
	PostCount      uint64 `json:"post_count"`
}

type Post struct {
	ID           uint64    `json:"ID"`
	ProfileID    uint64    `json:"profile_ID"`
	File         []byte    `json:"file"`
	Description  string    `json:"description"`
	LikeCount    uint64    `json:"like_count"`
	CommentCount uint64    `json:"comment_count"`
	DateTime     time.Time `json:"date_time"`
}

type Like struct {
	OwnerID   uint64 `json:"owner_ID"`
	OwnerName string `json:"owner_name"`
}

type Comment struct {
	ID            uint64    `json:"ID"`
	PostID        uint64    `json:"post_ID"`
	OwnerID       uint64    `json:"owner_ID"`
	OwnerUsername string    `json:"owner_username"`
	Text          string    `json:"text"`
	DateTime      time.Time `json:"date_time"`
}
