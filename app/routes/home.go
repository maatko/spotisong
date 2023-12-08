package routes

import (
	"net/http"
	"spotisong/api"

	"github.com/gorilla/mux"
)

type Home struct {
	Project api.ProjectInformation
}

func (home Home) Index(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusOK)

	home.Project = api.Project

	api.RenderTemplate(response, home, "base.html")
}

func (home Home) SetupRoutes(router *mux.Router) {
	router.HandleFunc("/", home.Index)
}
