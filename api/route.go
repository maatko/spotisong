package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type RouteCallback func(http.ResponseWriter, *http.Request)

type Route struct {
	Root string
	Callbacks map [string] RouteCallback
}

func (route Route) Register(router *mux.Router) {
	for path, callback := range route.Callbacks {
		router.HandleFunc(fmt.Sprintf("%s/%s", route.Root, path), callback)
	}
}