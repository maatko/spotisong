package routes

import (
	"fmt"
	"net/http"
	"spotisong/app/models"

	"github.com/gorilla/mux"
)

type Home struct {
}

func (home Home) Index(response http.ResponseWriter, request *http.Request) {
	user := models.User{ID: 1}

	err := user.Load("id")
	if err != nil {
		panic(err)
	}

	fmt.Println(user)
	// api.RenderTemplate(response, home, http.StatusOK, "base.html")
}

func (home Home) Login(response http.ResponseWriter, request *http.Request) {
	http.Redirect(response, request, "/auth/login", http.StatusPermanentRedirect)
}

func (home Home) Register(response http.ResponseWriter, request *http.Request) {
	http.Redirect(response, request, "/auth/register", http.StatusPermanentRedirect)
}

func (home Home) SetupRoutes(router *mux.Router) {
	router.HandleFunc("/", home.Index)
	router.HandleFunc("/login", home.Login)
	router.HandleFunc("/register", home.Register)
}
