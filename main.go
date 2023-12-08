package main

import (
	"fmt"
	"net/http"
	"os"
	"spotisong/api"
	"spotisong/app"
	"strings"

	"github.com/joho/godotenv"
)

func MakeMigrations(args []string) error {
	info, err := os.Stat(api.Project.MigrationsDirectory)
	if os.IsNotExist(err) {
		err = os.Mkdir(api.Project.MigrationsDirectory, 0755)
		if err != nil {
			return err
		}
	}

	if !info.IsDir() {
		return fmt.Errorf("'%v' must be a directory", api.Project.MigrationsDirectory)
	}

	fmt.Println("> Making migrations...")
	for _, migration := range api.Project.Migrations {
		file, err := os.Create(migration.GetFile())
		if err != nil {
			return err
		}
		defer file.Close()

		fmt.Printf("> table: %s\n", migration.Table)
		file.WriteString(migration.Schema)
	}

	return nil
}

func Migrate(args []string) error {
	info, err := os.Stat(api.Project.MigrationsDirectory)
	if os.IsNotExist(err) {
		return fmt.Errorf("'%v' does not exist, please run `makemigrations` first", api.Project.MigrationsDirectory)
	}

	if !info.IsDir() {
		return fmt.Errorf("'%v' must be a directory", api.Project.MigrationsDirectory)
	}

	// TODO :: expand the migration system to a more
	// advanced system where tables are altered not dropped,
	// and check the default values for columns.
	for _, migration := range api.Project.Migrations {
		migrationSchema, err := os.ReadFile(migration.GetFile())
		if err != nil {
			return err
		}

		currentSchema, err := migration.QuerySchema()
		if currentSchema != string(migrationSchema) || err != nil {
			fmt.Printf("Changes detected migrating '%v'\n", migration.Table)

			err = migration.Drop()
			if err != nil {
				return err
			}

			err = migration.Create()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func Watch(args []string) error {
	tailwind, err := api.NewTailWind(
		os.Getenv("TAILWIND_VERSION"),
		"./.tailwind",
	)

	if err != nil {
		return err
	}

	process, err := tailwind.Watch(
		api.Project.Src("style.css"),
		api.Project.Static("css/%s", os.Getenv("TAILWIND_OUTPUT")),
	)

	if err != nil {
		return err
	}

	fmt.Printf("[*] TailWindCSS exited [Exit Code: %v]\n", process.ExitCode())
	return nil
}

func Run(args []string) error {
	address := fmt.Sprintf(
		"%s:%s",
		os.Getenv("HTTP_ADDRESS"),
		os.Getenv("HTTP_PORT"),
	)

	api.Project.Router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./app/static/"))))

	fmt.Printf("> HTTP server listening on http://%s\n", address)
	return http.ListenAndServe(
		address,
		api.Project.Router,
	)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load '.env' file, maybe its missing?")
	}

	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println(ACTIONS_RESPONSE)
		return
	}

	err = api.Project.Setup(
		"./app",
		os.Getenv("APP_NAME"),
		os.Getenv("APP_STATIC_DIR"),
		"templates",
		os.Getenv("APP_MIGRATIONS_DIR"),
	)

	if err != nil {
		panic(err)
	}

	app.Initialize()

	if action, ok := ACTIONS[strings.ToLower(os.Args[1])]; ok {
		err = action(args[1:])
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println(ACTIONS_RESPONSE)
	}
}

var ACTIONS = map[string]func(args []string) error{
	"makemigrations": MakeMigrations,
	"migrate":        Migrate,
	"watch":          Watch,
	"run":            Run,
}

// this is the response that gets
// printed onto the screen if the
// user provided invalid launch args
const ACTIONS_RESPONSE = "<makemigrations/migrate/watch/run>"
