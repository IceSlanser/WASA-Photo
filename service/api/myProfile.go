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
	type UserName struct {
		Identifier string `json:"identifier"`
	}
	var userName UserName

	err := json.NewDecoder(r.Body).Decode(&userName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	isLegal, err := utils.IsLegal(userName.Identifier)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error during func isLegal")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !isLegal {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Inserisci nel database l'username, se esiste già
	// fai la query dell'id dell'username
	var id, exist = getId()
	if exist {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(201)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(id)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error during func getId")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

}

func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

}

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

}

func (rt *_router) setMyBiografy(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

}

func (rt *_router) setMyPropic(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

}
