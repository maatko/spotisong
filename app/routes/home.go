package routes

import (
	"html/template"
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

	tmpl := template.Must(template.ParseFiles("./app/templates/base.html"))
	tmpl.Execute(response, home)
}

func (home Home) SetupRoutes(router *mux.Router) {
	router.HandleFunc("/", home.Index)
}
