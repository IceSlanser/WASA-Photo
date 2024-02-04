package api

import (
	"encoding/json"
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !authorization {
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !authorization {
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Responses
	w.WriteHeader(http.StatusNoContent)
	w.Header().Set("Content-Type", "application/json")
}

func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

}

func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

}

func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

}

func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

}
