package api

import (
	"encoding/json"
	"github.com/IceSlanserUni/WASAPhoto/service/utils"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/IceSlanserUni/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// Get request's userProfile
func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get RequestUserID from the Header
	myUID, authorization, err := SecurityHandler(r, rt)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !authorization {
		ctx.Logger.WithError(err).Error("Operation not authorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userIDStr := ps.ByName("UID")
	userID, _ := strconv.Atoi(userIDStr)

	// Get DBUser
	dbProfile, err := rt.db.GetProfile(myUID, uint64(userID))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Convert DBUser into APIUser
	userProfile := NewUser(dbProfile)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(userProfile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Get request's userFullProfile
func (rt *_router) getFullProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get RequestUserID from the Header
	myUID, authorization, err := SecurityHandler(r, rt)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !authorization {
		ctx.Logger.WithError(err).Error("Operation not authorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get userProfile
	var userProfile User
	err = json.NewDecoder(r.Body).Decode(&userProfile)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get userPosts
	DBPosts, err := rt.db.GetPosts(myUID, userProfile.ID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// Convert DBPosts to APIPosts
	var APIPosts []Post
	for _, DBPost := range DBPosts {
		APIPost := NewPost(DBPost)
		APIPosts = append(APIPosts, APIPost)
	}

	// Get userFollows
	followings, followers, err := rt.db.GetFollows(myUID, userProfile.ID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(APIPosts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(followings)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(followers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Follow a certain user
func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get UserID from the Header
	UID, authorization, err := SecurityHandler(r, rt)
	if err != nil {
		ctx.Logger.WithError(err).Error("SecurityHandler has gone wrong")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !authorization {
		ctx.Logger.WithError(err).Error("Operation not authorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get followedUID from the router
	followedStr := ps.ByName("FollowedUID")
	followedUID, _ := strconv.Atoi(followedStr)

	// Put like
	followID, exist, err := rt.db.PutFollow(UID, uint64(followedUID))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if exist {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Responses
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(followID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Unfollow a certain user
func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get UserID from the Header
	UID, authorization, err := SecurityHandler(r, rt)
	if err != nil {
		ctx.Logger.WithError(err).Error("SecurityHandler has gone wrong")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !authorization {
		ctx.Logger.WithError(err).Error("Operation not authorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get followedUID from the router
	followedStr := ps.ByName("FollowedUID")
	followedUID, err := strconv.ParseUint(followedStr, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to parse request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Delete follow
	authorization, err = rt.db.DeleteFollow(UID, followedUID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if !authorization {
		ctx.Logger.WithError(err).Error("Operation not authorized")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Responses
	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get UserID from the Header
	UID, authorization, err := SecurityHandler(r, rt)
	if err != nil {
		ctx.Logger.WithError(err).Error("SecurityHandler has gone wrong")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !authorization {
		ctx.Logger.WithError(err).Error("Operation not authorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get bannedUID from the router
	bannedStr := ps.ByName("BannedUID")
	bannedUID, _ := strconv.Atoi(bannedStr)

	// Put like
	banID, exist, err := rt.db.PutBan(UID, uint64(bannedUID))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if exist {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Responses
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(banID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get UserID from the Header
	UID, authorization, err := SecurityHandler(r, rt)
	if err != nil {
		ctx.Logger.WithError(err).Error("SecurityHandler has gone wrong")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !authorization {
		ctx.Logger.WithError(err).Error("Operation not authorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get followedUID from the router
	bannedStr := ps.ByName("BannedUID")
	bannedUID, err := strconv.ParseUint(bannedStr, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to parse request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Delete follow
	authorization, err = rt.db.DeleteBan(UID, bannedUID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if !authorization {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Responses
	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get UserID from the Header
	UID, authorization, err := SecurityHandler(r, rt)
	if err != nil {
		ctx.Logger.WithError(err).Error("SecurityHandler has gone wrong")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !authorization {
		ctx.Logger.WithError(err).Error("Operation not authorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Limit From data to 10MB
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to parse form")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get the uploaded photoFile and description
	var photoFile []byte
	file, handler, err := r.FormFile("file")
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to get file from form")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()
	photoFile, err = ioutil.ReadAll(file)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to read form file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Check file extension
	if !utils.IsMediaFile(handler.Filename) {
		ctx.Logger.WithError(err).Error("Uploading a non media-type file")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the uploaded description
	description := r.FormValue("description")

	// Post photo
	postID, err := rt.db.PostPost(UID, photoFile, description)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to post a new photo")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Responses
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(postID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get UserID from the Header
	UID, authorization, err := SecurityHandler(r, rt)
	if err != nil {
		ctx.Logger.WithError(err).Error("SecurityHandler has gone wrong")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !authorization {
		ctx.Logger.WithError(err).Error("Operation not authorized")
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

	// Delete post
	authorization, err = rt.db.DeletePost(UID, postID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if !authorization {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Responses
	w.WriteHeader(http.StatusNoContent)
}
