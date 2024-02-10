package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	// MyProfile
	rt.router.POST("/session", rt.wrap(rt.doLogin))                  // V
	rt.router.GET("/stream", rt.wrap(rt.getMyStream))                // Controlla post
	rt.router.PUT("/profile/setUserName", rt.wrap(rt.setMyUserName)) // V

	// ManageProfile
	rt.router.GET("/users/:UID/profile", rt.wrap(rt.getUserProfile)) // Controlla dopo numero
	rt.router.PUT("/users/:UID/follow", rt.wrap(rt.followUser))
	rt.router.DELETE("/users/:UID/follow", rt.wrap(rt.unfollowUser))
	rt.router.PUT("/users/:UID/ban", rt.wrap(rt.banUser))
	rt.router.DELETE("/users/:UID/ban", rt.wrap(rt.unbanUser))
	rt.router.POST("/post", rt.wrap(rt.uploadPhoto))
	rt.router.DELETE("/post/:postID", rt.wrap(rt.deletePhoto))

	// ManagePost
	rt.router.POST("/post/:postID", rt.wrap(rt.getFullPost))
	rt.router.POST("/post/:postID/likes", rt.wrap(rt.likePhoto))
	rt.router.DELETE("/post/:postID/likes", rt.wrap(rt.unlikePhoto))
	rt.router.POST("/post/:postID/comments", rt.wrap(rt.commentPhoto))
	rt.router.DELETE("/post/:postID/comments/:commentUID", rt.wrap(rt.uncommentPhoto))

	return rt.router
}
