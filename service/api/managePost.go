package api

import (
	"encoding/json"
	"github.com/IceSlanserUni/WASAPhoto/service/database"
	"net/http"
	"strconv"

	"github.com/IceSlanserUni/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// getFullPost Get all related to a post
func (rt *_router) getFullPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get RequestUserID from the Header
	myUID, authorization, err := SecurityHandler(r, rt)
	if err != nil {
		ctx.Logger.WithError(err).Error("SecurityHandler has gone wrong")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !authorization {
		ctx.Logger.WithError(err).Error("getFullPost not authorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var FullPost struct {
		Post         Post          `json:"post"`
		LikeOwners   []string      `json:"like_owners"`
		FullComments []FullComment `json:"full_comments"`
	}

	// Get user's post
	postStr := ps.ByName("postID")
	postID, err := strconv.ParseUint(postStr, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to decode postID")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get post
	DBPost, err := rt.db.GetPostInfo(myUID, postID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to GetComments")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	FullPost.Post = NewPost(DBPost)

	// Get postComments
	var DBComments []database.Comment
	DBComments, err = rt.db.GetComments(myUID, postID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to GetComments")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// Convert DBPosts to APIComments
	for _, DBComment := range DBComments {
		var fullComment FullComment
		var temp database.User
		fullComment.Comment = NewComment(DBComment)
		temp, err = rt.db.GetProfile(myUID, fullComment.Comment.OwnerID)
		fullComment.Username = temp.Username
		FullPost.FullComments = append(FullPost.FullComments, fullComment)
	}

	// Get likeOwners
	FullPost.LikeOwners, err = rt.db.GetLikes(myUID, postID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to GetLikes")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(FullPost)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to encode FullPost")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// likePhoto put a like to a post
func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get UserID from the Header
	UID, authorization, err := SecurityHandler(r, rt)
	if err != nil {
		ctx.Logger.WithError(err).Error("SecurityHandler has gone wrong")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !authorization {
		ctx.Logger.WithError(err).Error("likePhoto not authorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get postID from the router
	postStr := ps.ByName("postID")
	postID, err := strconv.ParseUint(postStr, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to parse request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Put like
	exist, err := rt.db.PutLike(UID, postID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to PutLike")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if exist {
		ctx.Logger.WithError(err).Error("Post does not exist")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Responses
	w.WriteHeader(http.StatusNoContent)
}

// unlikePhoto delete a like from a post
func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get UserID from the Header
	UID, authorization, err := SecurityHandler(r, rt)
	if err != nil {
		ctx.Logger.WithError(err).Error("SecurityHandler has gone wrong")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !authorization {
		ctx.Logger.WithError(err).Error("unlikePhoto not authorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get postID from the router
	postStr := ps.ByName("postID")
	postID, err := strconv.ParseUint(postStr, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to parse request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Delete like
	authorization, err = rt.db.DeleteLike(UID, postID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to DeleteLike")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if !authorization {
		ctx.Logger.WithError(err).Error("DeleteLike not authorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Responses
	w.WriteHeader(http.StatusNoContent)
}

// commentPhoto comment a post
func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get UserID from the Header
	UID, authorization, err := SecurityHandler(r, rt)
	if err != nil {
		ctx.Logger.WithError(err).Error("SecurityHandler has gone wrong")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !authorization {
		ctx.Logger.WithError(err).Error("commentPhoto not authorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get postID from the router
	postStr := ps.ByName("postID")
	postID, err := strconv.ParseUint(postStr, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to parse request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the comment from the form
	text := r.FormValue("text")

	// Post Comment
	commentID, err := rt.db.PostComment(UID, postID, text)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to PostComment")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Responses
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	err = json.NewEncoder(w).Encode(commentID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to encode commentID")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

// uncommentPhoto delete a comment from a post
func (rt *_router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get UserID from the Header
	UID, authorization, err := SecurityHandler(r, rt)
	if err != nil {
		ctx.Logger.WithError(err).Error("SecurityHandler has gone wrong")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !authorization {
		ctx.Logger.WithError(err).Error("uncommentPhoto not authorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get postID from the router
	postStr := ps.ByName("postID")
	postID, err := strconv.ParseUint(postStr, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to parse request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get commentID from the router
	commentStr := ps.ByName("commentID")
	commentID, err := strconv.ParseUint(commentStr, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to parse request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Delete comment
	authorization, err = rt.db.DeleteComment(UID, postID, commentID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to DeleteComment")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if !authorization {
		ctx.Logger.WithError(err).Error("DeleteComment not authorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Responses
	w.WriteHeader(http.StatusNoContent)
}
