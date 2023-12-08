package routes

import (
	"net/http"
	"spotisong/api"

	"github.com/gorilla/mux"
)

type Auth struct {
}

func (auth Auth) Login(response http.ResponseWriter, request *http.Request) {
	api.RenderTemplate(
		response,
		auth,
		http.StatusOK,
		"base.html",
		"auth/base.html",
		"auth/login.html",
	)
}

func (auth Auth) Register(response http.ResponseWriter, request *http.Request) {
	api.RenderTemplate(
		response,
		auth,
		http.StatusOK,
		"base.html",
		"auth/base.html",
		"auth/register.html",
	)
}

func (auth Auth) SetupRoutes(router *mux.Router) {
	router.HandleFunc("/login", auth.Login)
	router.HandleFunc("/register", auth.Register)
}
