package routes

import (
	"net/http"
	"spotisong/api"

	"github.com/gorilla/mux"
)

type AuthenticationRoute struct {
	Route api.Route
}

func (authRoute AuthenticationRoute) Index(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Index Page"))
}

func (authRoute AuthenticationRoute) Register(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Register Page"))
}

func (authRoute AuthenticationRoute) Login(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Login Page"))
}

func (authRoute AuthenticationRoute) Logout(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Logout Page"))
}

func (auth AuthenticationRoute) Setup(router *mux.Router) {
	auth.Route.Callbacks = map [string] api.RouteCallback {
		"": auth.Index,
		"register/": auth.Register,
		"login/": auth.Login,
		"logout/": auth.Logout,
	}

	auth.Route.Register(router)
}