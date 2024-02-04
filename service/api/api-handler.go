package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	// MyProfile
	rt.router.POST("/session", rt.wrap(rt.doLogin))
	rt.router.GET("/stream", rt.wrap(rt.getMyStream))
	rt.router.GET("/users/:UID/profile", rt.wrap(rt.getUserProfile))
	rt.router.PUT("/profile/setUserName", rt.wrap(rt.setMyUserName))

	// ManageProfile
	rt.router.POST("/users/:UID/follow", rt.wrap(rt.followUser))
	rt.router.DELETE("/users/:UID/follow", rt.wrap(rt.unfollowUser))
	rt.router.POST("/users/:UID/ban", rt.wrap(rt.banUser))
	rt.router.DELETE("/users/:UID/ban", rt.wrap(rt.unbanUser))
	rt.router.POST("/post", rt.wrap(rt.uploadPhoto))
	rt.router.DELETE("/post/:postId", rt.wrap(rt.deletePhoto))

	// ManagePost
	rt.router.POST("/post/:postId/likes", rt.wrap(rt.likePhoto))
	rt.router.DELETE("/post/:postId/likes", rt.wrap(rt.unlikePhoto))
	rt.router.POST("/post/:postId/comments", rt.wrap(rt.commentPhoto))
	rt.router.DELETE("/post/:postId/comments/:commentuId", rt.wrap(rt.uncommentPhoto))

	return rt.router
}
