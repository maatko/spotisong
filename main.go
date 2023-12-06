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

type Callback func()
var Callbacks = map[string] Callback {
	"run": Run,
	"makemigrations": MakeMigrations,
	"migrate": Migrate,
}

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

	if callback, ok := Callbacks[strings.ToLower(args[0])]; ok {
		callback()
	} else {
		log.Fatal("Please choose an action: [run, makemigrations, migrate]")
	}
}

func Run() {
	fmt.Printf("HTTP server started on '%v:%v'...\n", api.Instance.Address, api.Instance.Port)
	http.ListenAndServe(fmt.Sprintf(
		"%v:%v",
		api.Instance.Address,
		api.Instance.Port,
	), nil)
}

func MakeMigrations() {
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

		// TODO :: in the future make this so it
		// doesn't drop the whole table but detect
		// the changes in the migration and if needed
		// it drops the table
		err = migration.Drop()
		if err != nil {
			panic(err)
		}
	}
}

func Migrate() {
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
}