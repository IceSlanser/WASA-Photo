package api

import (
	"time"

	"github.com/IceSlanserUni/WASAPhoto/service/database"
)

type User struct {
	ID        uint64 `json:"Id"`
	Username  string `json:"Username"`
	Following uint64 `json:"Following"`
	Followers uint64 `json:"Followers"`
	Post      uint64 `json:"Post"`
}

type Post struct {
	ID           uint64    `json:"ID"`
	ProfileID    string    `json:"profile_ID"`
	Description  string    `json:"description"`
	LikeCount    uint64    `json:"like_count"`
	CommentCount uint64    `json:"comment_count"`
	DateTime     time.Time `json:"date_time"`
}

func NewUser(user database.User) User {
	return User{
		ID:        user.ID,
		Username:  user.Username,
		Following: user.Following,
		Followers: user.Followers,
		Post:      user.Post,
	}
}

func NewPost(post database.Post) Post {
	return Post{
		ID:           post.ID,
		ProfileID:    post.ProfileID,
		Description:  post.Description,
		LikeCount:    post.LikeCount,
		CommentCount: post.CommentCount,
		DateTime:     post.DateTime,
	}
}
