package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/maatko/spotisong/models"

	"github.com/gorilla/mux"
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

	// make sure to update the database 
	// for the models
	models.DataBase = database

	/////////////////////////
	// Models
	/////////////////////////

	models.Register(models.User {
		Username: "256",
		Password: "512",
	})

	/////////////////////////

	callback(database)

	database.Close()
}

func MainHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)

	user := models.User {
		Username: "admin",
		Password: "password",
	}.Fetch()

	log.Printf("ID: %v, Username: %v, Password: %v\n", user.ID, user.Username, user.Password)
}

func Run(database *sql.DB) {
	log.Println("Starting HTTP server at port '8000'")

	router := mux.NewRouter()
    router.HandleFunc("/", MainHandler)

	http.Handle("/", router)
	http.ListenAndServe(":8000", nil)
}

func Migrate(database *sql.DB) {
	for _, model := range models.Models {
		var existingFields [] models.ModelField

		rows, err := database.Query(fmt.Sprintf("PRAGMA table_info(%v)", model.Table))
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			var cid int
			var name string
			var dataType string
			var notNull int
			var defaultValue sql.NullString
			var primaryKey int
	
			err := rows.Scan(&cid, &name, &dataType, &notNull, &defaultValue, &primaryKey)
			if err != nil {
				log.Fatal(err)
			}
			
			existingFields = append(existingFields, models.ModelField {
				Name: name,
				Type: dataType,
				Properties: "",
			})
		}

		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		if len(model.Fields) != len(existingFields) {
			database.Exec(model.GenerateMigrationSQL())

			log.Printf("Migrating the '%v' table...\n", model.Table)
			return
		}

		for index, field := range model.Fields {
			existingField := existingFields[index]
			if field.Name != existingField.Name || field.Type != existingField.Type {
				database.Exec(model.GenerateMigrationSQL())
	
				log.Printf("Migrating the '%v' table...\n", model.Table)
				return
			}
		}
	}
}