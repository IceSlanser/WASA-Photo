package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	// MyProfile
	rt.router.PUT("/session", rt.wrap(rt.doLogin))                   // V
	rt.router.GET("/stream", rt.wrap(rt.getMyStream))                // V
	rt.router.PUT("/profile/setUserName", rt.wrap(rt.setMyUserName)) // V

	// ManageProfile
	rt.router.GET("/users/:UID/profile", rt.wrap(rt.getUserProfile)) // V
	rt.router.PUT("/users/:UID/follow", rt.wrap(rt.followUser))      // V
	rt.router.DELETE("/users/:UID/follow", rt.wrap(rt.unfollowUser)) // V
	rt.router.PUT("/users/:UID/ban", rt.wrap(rt.banUser))            // V
	rt.router.DELETE("/users/:UID/ban", rt.wrap(rt.unbanUser))       // V
	rt.router.POST("/posts", rt.wrap(rt.uploadPhoto))                // V
	rt.router.DELETE("/posts/:postID", rt.wrap(rt.deletePhoto))      // V

	// ManagePost
	rt.router.GET("/posts/:postID", rt.wrap(rt.getFullPost))                           // V
	rt.router.PUT("/posts/:postID/likes", rt.wrap(rt.likePhoto))                       // V
	rt.router.DELETE("/posts/:postID/likes", rt.wrap(rt.unlikePhoto))                  // V
	rt.router.POST("/posts/:postID/comments", rt.wrap(rt.commentPhoto))                // V
	rt.router.DELETE("/posts/:postID/comments/:commentID", rt.wrap(rt.uncommentPhoto)) // V

	return rt.router
}
