package api

import (
	"time"

	"github.com/IceSlanserUni/WASAPhoto/service/database"
)

type User struct {
	ID             uint64 `json:"ID"`
	Username       string `json:"username"`
	FollowingCount uint64 `json:"following_count"`
	FollowersCount uint64 `json:"followers_count"`
	PostCount      uint64 `json:"post_count"`
}

type Post struct {
	ID           uint64    `json:"ID"`
	ProfileID    string    `json:"profile_ID"`
	File         []byte    `json:"file"`
	Description  string    `json:"description"`
	LikeCount    uint64    `json:"like_count"`
	CommentCount uint64    `json:"comment_count"`
	DateTime     time.Time `json:"date_time"`
}

type Comment struct {
	ID       uint64    `json:"ID"`
	PostID   uint64    `json:"post_ID"`
	OwnerID  uint64    `json:"owner_ID"`
	Text     string    `json:"text"`
	DateTime time.Time `json:"date_time"`
}

func NewUser(user database.User) User {
	return User{
		ID:             user.ID,
		Username:       user.Username,
		FollowingCount: user.FollowingCount,
		FollowersCount: user.FollowersCount,
		PostCount:      user.PostCount,
	}
}

func NewPost(post database.Post) Post {
	return Post{
		ID:           post.ID,
		ProfileID:    post.ProfileID,
		File:         post.File,
		Description:  post.Description,
		LikeCount:    post.LikeCount,
		CommentCount: post.CommentCount,
		DateTime:     post.DateTime,
	}
}

func NewComment(comment database.Comment) Comment {
	return Comment{
		ID:       comment.ID,
		PostID:   comment.PostID,
		OwnerID:  comment.OwnerID,
		Text:     comment.Text,
		DateTime: comment.DateTime,
	}
}
