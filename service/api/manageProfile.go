package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/IceSlanserUni/WASAPhoto/service/utils"

	"github.com/IceSlanserUni/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// getUserProfile Get request's userProfile
func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get RequestUserID from the Header
	myUID, authorization, err := SecurityHandler(r, rt)
	if err != nil {
		ctx.Logger.WithError(err).Error("SecurityHandler has gone wrong")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !authorization {
		ctx.Logger.WithError(err).Error("getUserProfile not authorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// userProfile struct
	var userProfile struct {
		Profile    User     `json:"profile"`
		Posts      []Post   `json:"posts"`
		Followings []uint64 `json:"followings"`
		Followers  []uint64 `json:"followers"`
	}

	// Get Profile
	userIDStr := ps.ByName("UID")
	userID, _ := strconv.Atoi(userIDStr)
	dbProfile, err := rt.db.GetProfile(myUID, uint64(userID))
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to GetProfile")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	userProfile.Profile = NewUser(dbProfile)

	// Get userPosts
	DBPosts, err := rt.db.GetUserPosts(myUID, userProfile.Profile.ID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to GetPosts")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// Convert DBPosts to APIPosts
	for _, DBPost := range DBPosts {
		APIPost := NewPost(DBPost)
		userProfile.Posts = append(userProfile.Posts, APIPost)
	}

	// Get userFollows
	userProfile.Followings, userProfile.Followers, err = rt.db.GetFollows(myUID, userProfile.Profile.ID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to encode GetFollows")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Responses
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(userProfile)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to encode userProfile")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// followUser Follow a certain user
func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get UserID from the Header
	UID, authorization, err := SecurityHandler(r, rt)
	if err != nil {
		ctx.Logger.WithError(err).Error("SecurityHandler has gone wrong")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !authorization {
		ctx.Logger.WithError(err).Error("followUser not authorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get followedUID from the router
	followedStr := ps.ByName("UID")
	followedUID, _ := strconv.Atoi(followedStr)

	// A user cannot follow himself
	if UID == uint64(followedUID) {
		ctx.Logger.WithError(err).Error("A user cannot follow himself")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Put like
	exist, err := rt.db.PutFollow(UID, uint64(followedUID))
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to PutFollow")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if !exist {
		ctx.Logger.WithError(err).Error("Followed does not exist")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Responses
	w.WriteHeader(http.StatusNoContent)

}

// unfollowUser Unfollow a certain user
func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get UserID from the Header
	UID, authorization, err := SecurityHandler(r, rt)
	if err != nil {
		ctx.Logger.WithError(err).Error("SecurityHandler has gone wrong")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !authorization {
		ctx.Logger.WithError(err).Error("unfollowUser not authorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get followedUID from the router
	followedStr := ps.ByName("UID")
	followedUID, err := strconv.ParseUint(followedStr, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to parse request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Delete follow
	authorization, err = rt.db.DeleteFollow(UID, followedUID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to DeleteFollow")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if !authorization {
		ctx.Logger.WithError(err).Error("DeleteFollow not authorized")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Responses
	w.WriteHeader(http.StatusNoContent)
}

// banUser ban a user
func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get UserID from the Header
	UID, authorization, err := SecurityHandler(r, rt)
	if err != nil {
		ctx.Logger.WithError(err).Error("SecurityHandler has gone wrong")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !authorization {
		ctx.Logger.WithError(err).Error("banUser not authorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get bannedUID from the router
	bannedStr := ps.ByName("UID")
	bannedUID, _ := strconv.Atoi(bannedStr)

	// A user cannot ban himself
	if UID == uint64(bannedUID) {
		ctx.Logger.WithError(err).Error("A user cannot ban himself")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Put like
	exist, err := rt.db.PutBan(UID, uint64(bannedUID))
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to PutBan")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if exist {
		ctx.Logger.WithError(err).Error("Banned does not exist")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Responses
	w.WriteHeader(http.StatusNoContent)
}

// unbanUser Unban a user
func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get UserID from the Header
	UID, authorization, err := SecurityHandler(r, rt)
	if err != nil {
		ctx.Logger.WithError(err).Error("SecurityHandler has gone wrong")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !authorization {
		ctx.Logger.WithError(err).Error("unbanUser not authorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get followedUID from the router
	bannedStr := ps.ByName("UID")
	bannedUID, err := strconv.ParseUint(bannedStr, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to parse request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Delete follow
	authorization, err = rt.db.DeleteBan(UID, bannedUID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to DeleteBan")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if !authorization {
		ctx.Logger.WithError(err).Error("DeleteBan not authorized")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Responses
	w.WriteHeader(http.StatusNoContent)
}

// uploadPhoto post a new Photo
func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get UserID from the Header
	UID, authorization, err := SecurityHandler(r, rt)
	if err != nil {
		ctx.Logger.WithError(err).Error("SecurityHandler has gone wrong")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !authorization {
		ctx.Logger.WithError(err).Error("uploadPhoto not authorized")
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Responses
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	err = json.NewEncoder(w).Encode(postID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to encode postID")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// deletePhoto delete a Photo
func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get UserID from the Header
	UID, authorization, err := SecurityHandler(r, rt)
	if err != nil {
		ctx.Logger.WithError(err).Error("SecurityHandler has gone wrong")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !authorization {
		ctx.Logger.WithError(err).Error("deletePhoto not authorized")
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
		ctx.Logger.WithError(err).Error("Failed to DeletePost")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if !authorization {
		ctx.Logger.WithError(err).Error("DeletePost not authorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Responses
	w.WriteHeader(http.StatusNoContent)
}
