package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Auth struct {
}

func (auth Auth) Login(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusOK)
	response.Write([]byte("Login page!"))
}

func (auth Auth) Register(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusOK)
	response.Write([]byte("Register page!"))
}

func (auth Auth) SetupRoutes(router *mux.Router) {
	router.HandleFunc("/login", auth.Login)
	router.HandleFunc("/register", auth.Register)
}
