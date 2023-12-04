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
	"github.com/joho/godotenv"
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

	err := godotenv.Load()
	if err != nil {
	  log.Fatal("Error loading .env file")
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

func Test() {
	rows, err := database.Select("user", "*", models.User {
		ID: 1,
		Username: "admin",
	})

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var id int
	var username string
	var password string

	for rows.Next() {
		err = rows.Scan(&id, &username, &password)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("ID:", id)
		log.Println("Username:", username)
		log.Println("Password:", password)
	}
}

func Run() {
	ip := os.Getenv("SERVER_IP")
	port := os.Getenv("SERVER_PORT")
	address := fmt.Sprintf("%v:%v", ip, port)

	log.Printf("Starting HTTP server at '%v'\n", address)

	router := mux.NewRouter()
	router.HandleFunc("/", MainHandler)

	http.ListenAndServe(address, router)
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