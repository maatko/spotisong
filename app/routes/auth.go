package routes

import (
	"fmt"
	"net/http"
	"spotisong/api"
	"spotisong/app/models"

	"github.com/gorilla/mux"
)

type Auth struct {
	Messages []api.Message
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
			api.ShowMessage("User does not exist!", true)
		} else {
			api.ShowMessage("User found!", false)
		}
	}

	api.RenderTemplate(
		response,
		map[string]any{
			"Messages": api.Project.Messages,
		},
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
	router.HandleFunc("/login", auth.Login)
	router.HandleFunc("/register", auth.Register)
}
