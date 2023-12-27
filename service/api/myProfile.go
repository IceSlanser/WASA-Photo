package api

import (
	"encoding/json"
	"net/http"

	"github.com/IceSlanserUni/WASAPhoto/service/api/reqcontext"
	"github.com/IceSlanserUni/WASAPhoto/service/utils"
	"github.com/julienschmidt/httprouter"
)

// doLogin checks if the user does exist.
// If the user does exist it'll return the HTTP Status 200 and return the user ID; otherwise, with HTTP Status 201.
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	username, err := utils.GetUsername(w, r, ctx)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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

	user, err, exist := rt.db.LoginUser(username)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error during func LoginUser")
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		if exist {
			w.WriteHeader(200)
		} else {
			w.WriteHeader(201)
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
	username, err := utils.GetUsername(w, r, ctx)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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

}

func (rt *_router) setMyBiografy(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

}

func (rt *_router) setMyPropic(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

}
