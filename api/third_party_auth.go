package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Engelbyte-s-Waterbyte/waterbyte-clean-architecture/models"
)

func FetchUserDataFromGoogle(token string) (user models.User, err error) {
	var res *http.Response
	res, err = http.Get("https://oauth2.googleapis.com/tokeninfo?id_token=" + token)
	if err != nil {
		return
	}
	var body map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		return
	}
	if body["error"] != nil {
		err = errors.New(body["error"].(string))
		return
	}

	user.Name = body["name"].(string)
	user.Email = body["email"].(string)
	gID := body["sub"].(string)
	user.GoogleID = &gID

	err = nil
	return
}

func FetchUserDataFromApple(token string) (user models.User, err error) {
	return
}
