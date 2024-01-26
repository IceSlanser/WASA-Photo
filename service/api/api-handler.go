package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	// MyProfile
	rt.router.POST("/session", rt.wrap(rt.doLogin))
	rt.router.GET("/", rt.wrap(rt.getMyStream))
	rt.router.GET("/users/:uId", rt.wrap(rt.getUserProfile))
	rt.router.PUT("/users/:uId", rt.wrap(rt.setMyUserName))

	// ManageProfile
	rt.router.PUT("/users/:uId/followers/:followerUId", rt.wrap(rt.followUser))
	rt.router.DELETE("/users/:uId/followers/:followerUId", rt.wrap(rt.unfollowUser))
	rt.router.POST("/users/:uId/banned/:banneduId", rt.wrap(rt.banUser))
	rt.router.DELETE("/users/:uId/banned/:banneduId", rt.wrap(rt.unbanUser))
	rt.router.POST("/users/:uId/post", rt.wrap(rt.uploadPhoto))
	rt.router.DELETE("/users/:uId/post/:postId", rt.wrap(rt.deletePhoto))

	// ManagePost
	rt.router.POST("/users/:uId/post/:postId/likes/:likeuId", rt.wrap(rt.likePhoto))
	rt.router.DELETE("/users/:uId/post/:postId/likes/:likeuId", rt.wrap(rt.unlikePhoto))
	rt.router.POST("/users/:uId/post/:postId/comments", rt.wrap(rt.commentPhoto))
	rt.router.DELETE("/users/:uId/post/:postId/comments/:commentuId", rt.wrap(rt.uncommentPhoto))

	return rt.router
}
