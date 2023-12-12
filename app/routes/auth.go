package routes

import (
	"fmt"
	"net/http"
	"spotisong/api"

	"github.com/gorilla/mux"
)

type Auth struct {
}

func (auth Auth) Login(response http.ResponseWriter, request *http.Request) {
	// authentication, _ := api.AppCookieStore.Get(request, "authentication")

	if request.Method == "POST" {
		fmt.Println("Login / POST")
	}

	api.RenderRoute(response, "auth", "login.html", auth)
}

func (auth Auth) Register(response http.ResponseWriter, request *http.Request) {
	// authentication, _ := api.AppCookieStore.Get(request, "authentication")

	if request.Method == "POST" {
		fmt.Println("Register / POST")
	}

	api.RenderRoute(response, "auth", "register.html", auth)
}

func (auth Auth) SetupRoutes(router *mux.Router) {
	router.HandleFunc("/login", auth.Login)
	router.HandleFunc("/register/", auth.Register)
}

var Authentication Auth = Auth{}
