package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/IceSlanserUni/WASAPhoto/service/api/reqcontext"
	"github.com/IceSlanserUni/WASAPhoto/service/utils"
	"github.com/julienschmidt/httprouter"
)

// doLogin checks if the user does exist.
// If the user does exist it'll return the HTTP Status 200 and return the user ID; otherwise, with HTTP Status 201.
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the username from the requestBody
	username, err := utils.GetUsername(r)
	if err != nil {
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
	user, exist, err := rt.db.LoginUser(username)
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
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(user.UserId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

}

func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

}

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the UserId from the Header
	uId, authorization, err := SecurityHandler(r, rt)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !authorization {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get the new username from the requestBody
	nname, err := utils.GetUsername(r)
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
	err = rt.db.SetUsername(uId, nname)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error during func SetUsername")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(uId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (rt *_router) setMyBiografy(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

}

func (rt *_router) setMyPropic(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

}

func SecurityHandler(r *http.Request, rt *_router) (uint64, bool, error) {
	authHeader := r.Header.Get("Authorization")

	available, err := rt.db.IsAvailable(authHeader)
	if err != nil {
		return 0, false, err
	}
	if available {
		return 0, false, nil
	}
	res, _ := strconv.Atoi(authHeader)

	return uint64(res), true, nil
}
