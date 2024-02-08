package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/IceSlanserUni/WASAPhoto/service/api/reqcontext"
	"github.com/IceSlanserUni/WASAPhoto/service/utils"
	"github.com/julienschmidt/httprouter"
)

// If the username does not exist, it will create a new profile, and an identifier is returned.
// If the username exists, the profile identifier is returned.
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the username from the requestBody
	username, err := utils.GetMyUsername(r)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to parse request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the input username is legal
	isLegal, err := utils.IsLegal(username)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error during func isLegal")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !isLegal {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// User Login
	userID, exist, err := rt.db.LoginUser(username)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error during func LoginUser")
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		if exist {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusCreated)
		}
	}

	// Responses
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// The stream is composed by photos from “following” (other users that the user follows).
func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	var TimeRange struct {
		StartTime time.Time `json:"start_time"`
		EndTime   time.Time `json:"end_time"`
	}

	err = json.NewDecoder(r.Body).Decode(&TimeRange)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to parse request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get list of DBPosts
	DBPosts, err := rt.db.GetStream(UID, TimeRange.StartTime, TimeRange.EndTime)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error during func GetPosts")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Convert DBPosts to APIPosts
	var APIPosts []Post
	for _, DBPost := range DBPosts {
		APIPost := NewPost(DBPost)
		APIPosts = append(APIPosts, APIPost)
	}

	// Responses
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(APIPosts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

// Set a new username for the current user.
func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	// Get the new username from the requestBody
	nname, err := utils.GetMyUsername(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the input username is legal
	isLegal, err := utils.IsLegal(nname)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error during func IsLegal")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !isLegal {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the input username is already taken
	available, err := rt.db.IsAvailable(nname)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error during func IsAvailable")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !available {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Update the Username
	err = rt.db.SetUsername(UID, nname)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error during func SetUsername")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(UID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func SecurityHandler(r *http.Request, rt *_router) (uint64, bool, error) {
	authentication := r.Header.Get("Authorization")

	available, err := rt.db.IsAvailable(authentication)
	if err != nil {
		return 0, false, err
	}
	if available {
		return 0, false, nil
	}
	res, _ := strconv.Atoi(authentication)

	return uint64(res), true, nil
}
