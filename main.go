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
		log.Fatal("Please choose an action: [run, makemigrations, migrate]")
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
		fmt.Printf("HTTP server started on '%v:%v'...\n", api.Instance.Address, api.Instance.Port)
		http.ListenAndServe(fmt.Sprintf(
			"%v:%v",
			api.Instance.Address,
			api.Instance.Port,
		), nil)
	} else if action == "migrate" {
		files, err := os.ReadDir("./app/migrations")
		if err != nil {
			panic(err)
		}

		for _, file := range files {
			bytes, err := os.ReadFile(fmt.Sprintf("%s/%s", "./app/migrations", file.Name()))
			if err != nil {
				panic(err)
			}

			_, err = api.Instance.DataBase.Exec(string(bytes))
			if err != nil {
				panic(err)
			}
		}
	} else if action == "makemigrations" {
		for _, model := range api.ModelRegistry {
			api.CreateMigration(&model)
		}

		for _, migration := range api.MigrationRegistry {
			file, err := os.Create(fmt.Sprintf("%v/%v-%v.sql", "./app/migrations", migration.Index, migration.Table))
			if err != nil {
				panic(err)
			}
			defer file.Close()

			_, err = file.WriteString(migration.Query)
			if err != nil {
				panic(err)
			}

			err = migration.Drop()
			if err != nil {
				panic(err)
			}
		}
	} else {
		log.Fatal("Please choose an action: [run, makemigrations, migrate]")
	}
}