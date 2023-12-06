package main

import (
	"fmt"
	"io"
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
	"watch": Watch,
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
	
	if callback, ok := Callbacks[strings.ToLower(args[0])]; ok {
		callback()
	} else {
		log.Fatal("Please choose an action: [run, makemigrations, migrate]")
	}
}

func Run() {
	handler := app.OnRouteRegister()
	
	log.Printf("HTTP server started on '%v:%v'...\n", api.Instance.Address, api.Instance.Port)
	http.ListenAndServe(fmt.Sprintf(
		"%v:%v",
		api.Instance.Address,
		api.Instance.Port,
	), handler)
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

func Watch() {
	_, err := os.Stat("./.tailwind")
	if os.IsNotExist(err) {
		err = os.Mkdir(".tailwind", os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	_, err = os.Stat("./.tailwind/tailwind")
	if os.IsNotExist(err) {
		log.Println("Failed to locate tailwind, downloading the binary...")

		response, err := http.Get("https://github.com/tailwindlabs/tailwindcss/releases/download/v3.3.6/tailwindcss-linux-x64")
		if err != nil {
			panic(err)
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			log.Fatalf("Error: Unexpected status code [%v]\n", response.Status)
		}

		destFile, err := os.Create("./.tailwind/tailwind")
		if err != nil {
			panic(err)
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, response.Body)
		if err != nil {
			panic(err)
		}
		
		log.Println("TailwindCSS binary downloaded!")

		err = destFile.Chmod(0755)
		if err != nil {
			panic(err)
		}
	}

	process, err := os.StartProcess(
		"./.tailwind/tailwind",
		[] string {
			"./.tailwind/tailwind",
			"-i", "./app/style.css",
			"-o", "./app/static/global.css",
			"--watch",
		},
		&os.ProcAttr {
			Files: [] *os.File {
				os.Stdin, 
				os.Stdout, 
				os.Stderr,
			},
		},
	)

	if err != nil {
		panic(err)
	}

	state, err := process.Wait()
	if err != nil {
		panic(err)
	}

	fmt.Println("Process finished with exit code:", state.ExitCode())
}