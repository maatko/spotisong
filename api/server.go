package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type HttpRoute interface {
	SetupRoutes(*mux.Router)
}

type HttpRoutes map[string]HttpRoute

type HttpServer struct {
	DataBase *sql.DB
	Router   *mux.Router
}

var Server HttpServer = HttpServer{}

func (server *HttpServer) Initialize(registerRoutes func() HttpRoutes) error {
	dataBase, err := sql.Open("sqlite3", os.Getenv("DB_CONNECTION"))
	if err != nil {
		return err
	}

	server.DataBase = dataBase

	server.Router = mux.NewRouter()
	server.Router.StrictSlash(true)

	// this registers a route for
	// all the resources so they can be accessed
	resourcesRouter := server.Router.PathPrefix("/resource/")
	resourcesRouter.Handler(http.StripPrefix(
		"/resource/",
		http.FileServer(http.Dir(GetResource("")))))

	for path, route := range registerRoutes() {
		router := server.Router.PathPrefix(path).Subrouter()
		route.SetupRoutes(router)
	}

	return nil
}

func (server *HttpServer) Start() error {
	address := os.Getenv("APP_URL")
	fmt.Println("HTTP server listening at", address)

	address = strings.Split(address, "://")[1]
	return http.ListenAndServe(address, server.Router)
}
