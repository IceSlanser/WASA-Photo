package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/IceSlanserUni/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

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
	followedUID, err := strconv.ParseUint(followedStr, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to parse request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Put like
	followID, exist, err := rt.db.PutFollow(UID, followedUID)
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
