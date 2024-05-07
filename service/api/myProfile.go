package api

import (
	"encoding/json"
	"net/http"

	"github.com/IceSlanserUni/WASAPhoto/service/api/reqcontext"
	"github.com/IceSlanserUni/WASAPhoto/service/database"
	"github.com/julienschmidt/httprouter"
)

//	doLogin If the username does not exist, it will create a new profile, and an identifier is returned.
//
// If the username exists, the profile identifier is returned.
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the username from the requestBody
	username, err := GetMyUsername(r)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to parse request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the input username is legal
	isLegal, err := IsLegal(username)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error during func isLegal")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !isLegal {
		ctx.Logger.WithError(err).Error("UserName is illegal")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// User Login
	DBUser, exist, err := rt.db.LoginUser(username)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error during func LoginUser")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Responses
	if exist {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(DBUser.ID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to encode userID")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// getMyStream The stream is composed by photos from “following” (other users that the user follows).
func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get UserID from the Header
	UID, authorization, err := SecurityHandler(r, rt)
	if err != nil {
		ctx.Logger.WithError(err).Error("SecurityHandler has gone wrong")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !authorization {
		ctx.Logger.WithError(err).Error("getMyStream not authorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get list of DBPosts
	DBPosts, err := rt.db.GetStream(UID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error during func GetStream")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get list of fullPosts
	var stream []FullPost
	for _, DBPost := range DBPosts {
		var fullPost FullPost
		fullPost.Post = NewPost(DBPost)
		fullPost.Post.DateTime = DBPost.DateTime.Format("2006-01-02 15:04:05")
		fullPost.Post.Username, err = rt.IDtoUsername(DBPost.ProfileID)
		if err != nil {
			ctx.Logger.WithError(err).Error("Failed to IDtoUsername")
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Get postComments
		var DBComments []database.Comment
		DBComments, err = rt.db.GetComments(UID, fullPost.Post.ID)
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
			temp, err = rt.db.GetProfile(UID, fullComment.Comment.OwnerID)
			if err != nil {
				ctx.Logger.WithError(err).Error("Failed to GetProfile")
				w.WriteHeader(http.StatusNotFound)
				return
			}
			fullComment.Username = temp.Username
			fullComment.Comment.DateTime = DBComment.DateTime.Format("2006-01-02 15:04:05")
			fullPost.FullComments = append(fullPost.FullComments, fullComment)
		}

		// Get likeOwners
		var DBLikes []database.Like
		DBLikes, err = rt.db.GetLikes(UID, fullPost.Post.ID)
		if err != nil {
			ctx.Logger.WithError(err).Error("Failed to GetLikes")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		// Convert DBPosts to APIComments
		for _, DBLike := range DBLikes {
			var like Like
			like = NewLike(DBLike)
			fullPost.LikeOwners = append(fullPost.LikeOwners, like)
		}

		stream = append(stream, fullPost)
	}

	// Responses
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(stream)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to encode stream")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

// setMyUserName Set a new username for the current user.
func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get UserID from the Header
	UID, authorization, err := SecurityHandler(r, rt)
	if err != nil {
		ctx.Logger.WithError(err).Error("SecurityHandler has gone wrong")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !authorization {
		ctx.Logger.WithError(err).Error("setMyUserName not authorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get the new username from the requestBody
	newName, err := GetMyUsername(r)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to GetMyUserName")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the input username is legal
	isLegal, err := IsLegal(newName)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error during func IsLegal")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !isLegal {
		ctx.Logger.WithError(err).Error("newName is illegal")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the input username is already taken
	_, available := rt.db.IsAvailable(newName)
	if !available {
		ctx.Logger.WithError(err).Error("newName is not available")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Update the Username
	err = rt.db.SetUsername(UID, newName)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error during func SetUsername")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// SecurityHandler Get UID from the Header
func SecurityHandler(r *http.Request, rt *_router) (uint64, bool, error) {
	authentication := r.Header.Get("Authorization")

	UID, available := rt.db.IsAvailable(authentication)
	if available {
		return 0, false, nil
	}

	return UID, true, nil
}
