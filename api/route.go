package api

import "github.com/gorilla/mux"

type RouteHandler interface {
	SetupRoutes(*mux.Router)
}
