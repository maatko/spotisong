package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/maatko/spotisong/database"
	"github.com/maatko/spotisong/models"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type Callback func()
var Callbacks = map[string] Callback {
	"run": Run,
	"migrate": Migrate,
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatal("Please choose an action: [run, migrate]")
	}

	database.Initialize("db.sqlite3")

	/////////////////////////
	// Models
	/////////////////////////

	models.Register(models.User {
		Username: "256",
		Password: "512",
	})

	/////////////////////////
	// Processing
	/////////////////////////
	
	if callback, ok := Callbacks[strings.ToLower(args[0])]; ok {
		callback()
	} else {
		log.Fatal("Please choose an action: [run, migrate]")
	}

	database.Close()
}

func MainHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)

	user := models.User {
		Username: "admin",
		Password: "password",
	}.Fetch()

	w.Write([]byte(fmt.Sprintf("ID: %v, Username: %v, Password: %v", user.ID, user.Username, user.Password)))
}

func Run() {
	log.Println("Starting HTTP server at port '8000'")

	router := mux.NewRouter()
	router.HandleFunc("/", MainHandler)

	http.ListenAndServe(":8000", router)
}

func Migrate() {
	for _, model := range models.Models {
		columns, err := database.GetTableInformation(model.Table)
		if err != nil {
			log.Fatal(err)
		}

		if len(model.Fields) != len(columns) {
			model.Migrate()
			return
		}

		for index, field := range model.Fields {
			existingField := columns[index]
			if field.Name != existingField.Name || field.Type != existingField.Type {
				model.Migrate()
				return
			}
		}
	}
}