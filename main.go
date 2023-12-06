package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"spotisong/api"
	"spotisong/app"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatal("Please choose an action: [run, migrate]")
	}

	err := godotenv.Load()
	if err != nil {
	  log.Fatal("Error loading .env file")
	}

	// make sure that the port of the server
	// is a valid integer
	port, err := strconv.Atoi(os.Getenv("HTTP_PORT"))
	if err != nil {
		log.Fatal("Port of the HTTP server must be a valid number")
	}

	// api has to be initialized before any migrations
	// or models are created. Becuase it has the main
	// instance of the database
	api.Instance = api.API {
		Address: os.Getenv("HTTP_ADDRESS"),
		Port: port,
	}.Initialize(os.Getenv("DB_CONNECTION"))

	// register all models & routes in the
	// main application
	app.OnModelRegister()
	app.OnRouteRegister()

	action := strings.ToLower(args[0])
	if action == "run" {
		// TODO :: run the http server

		fmt.Printf("HTTP server started on '%v:%v'...\n", api.Instance.Address, api.Instance.Port)
		http.ListenAndServe(fmt.Sprintf(
			"%v:%v",
			api.Instance.Address,
			api.Instance.Port,
		), nil)
	} else if action == "migrate" {
		// TODO :: run all the migrations
	} else {
		log.Fatal("Please choose an action: [run, migrate]")
	}
}