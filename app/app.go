package app

import (
	"net/http"
	"spotisong/api"
	"spotisong/app/models"
	"spotisong/app/routes"

	"github.com/gorilla/mux"
)

func OnRouteRegister() http.Handler {
	/////////////////////////////////////////
	// Register all your routes here
	/////////////////////////////////////////
	router := mux.NewRouter().StrictSlash(true)

	routes.AuthenticationRoute {
		Route: api.Route {
			Root: "",
		},
	}.Setup(router)

	return router
}

func OnModelRegister() {
	/////////////////////////////////////////
	// Register all your models here
	/////////////////////////////////////////

	models.User {}.Register()
	models.Post {}.Register()
}