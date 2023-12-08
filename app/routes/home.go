package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Home struct {
}

func (home Home) Index(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusOK)
	response.Write([]byte("Index page!"))
}

func (home Home) SetupRoutes(router *mux.Router) {
	router.HandleFunc("/", home.Index)
}
