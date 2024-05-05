package api

import (
	"github.com/IceSlanserUni/WASAPhoto/service/database"
	"strconv"
)

type User struct {
	ID             uint64 `json:"ID"`
	Username       string `json:"username"`
	FollowingCount uint64 `json:"following_count"`
	FollowerCount  uint64 `json:"follower_count"`
	PostCount      uint64 `json:"post_count"`
}

type Post struct {
	ID           uint64 `json:"ID"`
	Username     string `json:"username"`
	File         []byte `json:"file"`
	Description  string `json:"description"`
	LikeCount    uint64 `json:"like_count"`
	CommentCount uint64 `json:"comment_count"`
	DateTime     string `json:"date_time"`
}

type Comment struct {
	ID       uint64 `json:"ID"`
	PostID   uint64 `json:"post_ID"`
	OwnerID  uint64 `json:"owner_ID"`
	Text     string `json:"text"`
	DateTime string `json:"date_time"`
}

type FullComment struct {
	Username string  `json:"username"`
	Comment  Comment `json:"comment"`
}

type FullPost struct {
	Post         Post          `json:"post"`
	LikeOwners   []string      `json:"like_owners"`
	FullComments []FullComment `json:"full_comments"`
}

func NewUser(user database.User) User {
	return User{
		ID:             user.ID,
		Username:       user.Username,
		FollowingCount: user.FollowingCount,
		FollowerCount:  user.FollowerCount,
		PostCount:      user.PostCount,
	}
}

func NewPost(post database.Post) Post {
	return Post{
		ID:           post.ID,
		Username:     strconv.FormatUint(post.ProfileID, 10),
		File:         post.File,
		Description:  post.Description,
		LikeCount:    post.LikeCount,
		CommentCount: post.CommentCount,
		DateTime:     "",
	}
}

func NewComment(comment database.Comment) Comment {
	return Comment{
		ID:       comment.ID,
		PostID:   comment.PostID,
		OwnerID:  comment.OwnerID,
		Text:     comment.Text,
		DateTime: "",
	}
}
