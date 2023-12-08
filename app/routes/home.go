package routes

import (
	"net/http"
	"spotisong/api"

	"github.com/gorilla/mux"
)

type Home struct {
}

func (home Home) Index(response http.ResponseWriter, request *http.Request) {
	api.RenderTemplate(response, home, http.StatusOK, "base.html")
}

func (home Home) SetupRoutes(router *mux.Router) {
	router.HandleFunc("/", home.Index)
}
