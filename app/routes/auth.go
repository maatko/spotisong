package routes

import (
	"fmt"
	"net/http"
	"spotisong/api"
	"spotisong/app/models"

	"github.com/gorilla/mux"
)

type Auth struct {
}

func (auth Auth) Index(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusOK)
	response.Write([]byte("Index"))
}

func (auth Auth) Login(response http.ResponseWriter, request *http.Request) {
	var status int = http.StatusOK

	if request.Method == "POST" {
		user := models.User{}

		err := user.FromRequest(request)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = user.Fetch("username", "password")
		if err != nil {
			status = http.StatusUnauthorized
		} else {
			status = http.StatusOK
		}
	}

	api.RenderTemplate(
		response,
		auth,
		status,
		"base.html",
		"auth/base.html",
		"auth/login.html",
	)
}

func (auth Auth) Register(response http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		fmt.Println("Register / POST")
	}

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
	router.HandleFunc("/", auth.Index)
	router.HandleFunc("/login", auth.Login)
	router.HandleFunc("/register/", auth.Register)
}
