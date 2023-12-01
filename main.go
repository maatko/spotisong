package main

import (
	"database/sql"
	"log"
	"os"
	"strings"

	"github.com/maatko/spotisong/migrations"
	"github.com/maatko/spotisong/models"
	_ "github.com/mattn/go-sqlite3"
)

type Callback func(db *sql.DB)
var Callbacks = map[string] Callback {
	"run": Run,
	"migrate": Migrate,
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatal("Please choose an action: [run, migrate]")
	}

	callback := Callbacks[strings.ToLower(args[0])]
	if callback == nil {
		log.Fatal("Please choose an action: [run, migrate]")
	}

	database, err := sql.Open("sqlite3", "db.sqlite3")
	if err != nil {
		log.Fatal("Failed to open connection to the database")
	}
	defer database.Close()

	callback(database)
}

func Run(database *sql.DB) {
}

func Migrate(database *sql.DB) {
	// create all the migrations here
	migrations.Create("user", models.User {
		Email: "1024",
		Password: "512",
	})

	// run all the migrations here
	migrations.Migrate(database)
}