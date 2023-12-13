package routes

import (
	"net/http"
	"spotisong/api"
	"spotisong/app/models"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

const AppSessionCookie = "session_id"

type Auth struct {
}

func (auth Auth) Login(response http.ResponseWriter, request *http.Request) {
	_, err := models.GetCookieSession(request)
	if err == nil {
		http.Redirect(response, request, "/", http.StatusTemporaryRedirect)
		return
	}

	if request.Method == "POST" {
		user, err := models.NewRequestUser(request, false, false)
		if err != nil {
			api.MessageError(err.Error())
			RenderRoute(response, request, "auth", "login.html", auth)
			return
		}

		// used for checking against the cache
		password := user.Password

		err = user.Load("username")
		if err != nil || user.ID == 0 {
			api.MessageError("Invalid credentials!")
			RenderRoute(response, request, "auth", "login.html", auth)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			api.MessageError("Invalid credentials!")
			RenderRoute(response, request, "auth", "login.html", auth)
			return
		}

		_, err = models.NewCookieSession(request, response, *user)
		if err != nil {
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}

		api.MessageInfo("You have logged in!")
		http.Redirect(response, request, "/", http.StatusTemporaryRedirect)
		return
	}

	RenderRoute(response, request, "auth", "login.html", auth)
}

func (auth Auth) Logout(response http.ResponseWriter, request *http.Request) {
	// todo :: log out of the current session
}

func (auth Auth) Register(response http.ResponseWriter, request *http.Request) {
	_, err := models.GetCookieSession(request)
	if err == nil {
		http.Redirect(response, request, "/", http.StatusTemporaryRedirect)
		return
	}

	if request.Method == "POST" {
		user, err := models.NewRequestUser(request, true, true)
		if err != nil {
			api.MessageError(err.Error())
			RenderRoute(response, request, "auth", "register.html", auth)
			return
		}

		err = user.Save()
		if err != nil {
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = models.NewCookieSession(request, response, *user)
		if err != nil {
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}

		api.MessageInfo("You have been registered!")
		http.Redirect(response, request, "/", http.StatusTemporaryRedirect)
		return
	}

	RenderRoute(response, request, "auth", "register.html", auth)
}

func (auth Auth) SetupRoutes(router *mux.Router) {
	router.HandleFunc("/login", auth.Login).Methods("GET", "POST")
	router.HandleFunc("/logout", auth.Logout).Methods("GET", "POST")
	router.HandleFunc("/register", auth.Register).Methods("GET", "POST")
}

var Authentication Auth = Auth{}
