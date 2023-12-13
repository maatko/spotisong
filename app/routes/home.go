package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Home struct {
}

func (home Home) Index(response http.ResponseWriter, request *http.Request) {
	RenderRoute(response, request, "home", "index.html", home)
}

func (home Home) Login(response http.ResponseWriter, request *http.Request) {
	http.Redirect(response, request, "/auth/login", http.StatusPermanentRedirect)
}

func (home Home) Logout(response http.ResponseWriter, request *http.Request) {
	http.Redirect(response, request, "/auth/logout", http.StatusPermanentRedirect)
}

func (home Home) Register(response http.ResponseWriter, request *http.Request) {
	http.Redirect(response, request, "/auth/register", http.StatusPermanentRedirect)
}

func (home Home) SetupRoutes(router *mux.Router) {
	router.HandleFunc("/", home.Index)
	router.HandleFunc("/login", home.Login)
	router.HandleFunc("/logout", home.Logout)
	router.HandleFunc("/register", home.Register)
}
