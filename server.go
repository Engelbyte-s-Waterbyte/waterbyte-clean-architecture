package main

import (
	"log"
	"net/http"

	"github.com/Engelbyte-s-Waterbyte/waterbyte-clean-architecture/api"
	"github.com/Engelbyte-s-Waterbyte/waterbyte-clean-architecture/db"
	"github.com/Engelbyte-s-Waterbyte/waterbyte-clean-architecture/routes"
	"github.com/julienschmidt/httprouter"
)

func main() {
	server := httprouter.New()

	server.POST("/authentication/sign-in", routes.SignInRoute(db.SelectUserByThirdPartyID, api.FetchUserDataFromGoogle, api.FetchUserDataFromApple, db.SelectNextUsername, db.InsertUser))
	server.GET("/authentication/authenticated-user", routes.AuthenticatedUserRoute(db.SelectUserByID))

	log.Fatal("Error serving:", http.ListenAndServe(":9999", server))
}
