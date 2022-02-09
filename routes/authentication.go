package routes

import (
	"encoding/json"
	"net/http"

	"github.com/Engelbyte-s-Waterbyte/waterbyte-clean-architecture/logic"
	"github.com/Engelbyte-s-Waterbyte/waterbyte-clean-architecture/models"
	"github.com/julienschmidt/httprouter"
)

type signInRequest struct {
	GoogleToken *string `json:"google_token"`
	AppleToken  *string `json:"apple_token"`
}

func SignInRoute(selectUserByThirdPartyID func(*models.User) (found bool, err error), fetchUserDataFromGoogle, fetchUserDataFromApple logic.FetchUserDataFunc, selectNextUseraname logic.SelectNextUsernameFunc, insertUser func(*models.User) error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		jsonResponse(w)

		var req signInRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			badRequest(w)
			return
		}

		token, err := logic.SignIn(req.GoogleToken, req.AppleToken, selectUserByThirdPartyID, fetchUserDataFromGoogle, fetchUserDataFromApple, selectNextUseraname, insertUser)
		if err != nil {
			errResponse(w, err.Error())
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{"token": token})
	}
}

func AuthenticatedUserRoute(selectUserByID func(*models.User) (found bool, err error)) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		jsonResponse(w)

		token := r.Header.Get("Authorization")
		user, err := logic.AuthenticatedUser(token, selectUserByID)
		if err != nil {
			errResponse(w, err.Error())
			return
		}

		json.NewEncoder(w).Encode(user)
	}
}
