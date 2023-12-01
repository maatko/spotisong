package main

import (
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
	"test": Test,
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
	}.Create()

	log.Printf("ID: %v, Username: %v, Password: %v\n", user.ID, user.Username, user.Password)
}

func Test() {
	rows, err := database.GetTableInformation("user")
	if err != nil {
		log.Fatal(err)
	}

	for _, information := range rows {
		log.Println("======")
		log.Printf("ID: %v\n", information.ID)
		log.Printf("Name: %v\n", information.Name)
		log.Printf("Type: %v\n", information.Type)
		log.Printf("NonNull: %v\n", information.NonNull)
		log.Printf("Default Value String: %v\n", information.DefaultValue.String)
		log.Printf("Default Value Valid: %v\n", information.DefaultValue.Valid)
		log.Printf("PrimaryKey: %v\n", information.PrimaryKey)
	}
}


func Run() {
	log.Println("Starting HTTP server at port '8000'")

	router := mux.NewRouter()

    router.HandleFunc("/", MainHandler)

	http.Handle("/", router)
	
	http.ListenAndServe(":8000", nil)
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