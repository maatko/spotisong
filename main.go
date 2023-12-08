package main

import (
	"database/sql"
	"fmt"
	"os"
	"spotisong/api"
	"spotisong/app"
	"strings"

	"github.com/joho/godotenv"
)

func MakeMigrations(args []string) error {
	info, err := os.Stat(api.MIGRATIONS_DIRECTORY)
	if os.IsNotExist(err) {
		err = os.Mkdir(api.MIGRATIONS_DIRECTORY, 0755)
		if err != nil {
			return err
		}
	}

	if !info.IsDir() {
		return fmt.Errorf("'%v' must be a directory", api.MIGRATIONS_DIRECTORY)
	}

	fmt.Println("> Making migrations...")
	for _, migration := range api.Migrations {
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
	info, err := os.Stat(api.MIGRATIONS_DIRECTORY)
	if os.IsNotExist(err) {
		return fmt.Errorf("'%v' does not exist, please run `makemigrations` first", api.MIGRATIONS_DIRECTORY)
	}

	if !info.IsDir() {
		return fmt.Errorf("'%v' must be a directory", api.MIGRATIONS_DIRECTORY)
	}

	// TODO :: expand the migration system to a more
	// advanced system where tables are altered not dropped,
	// and check the default values for columns.
	for _, migration := range api.Migrations {
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
	tailwind := api.TailWind{
		Version: os.Getenv("TAILWIND_VERSION"),
		Binary:  "./.tailwind/",
	}

	return tailwind.Watch(
		"./app/style.css",
		"./app/static/"+os.Getenv("TAILWIND_OUTPUT"),
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

	api.DataBase, err = sql.Open("sqlite3", os.Getenv("DATABASE_CONNECTION"))
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
}

// this is the response that gets
// printed onto the screen if the
// user provided invalid launch args
const ACTIONS_RESPONSE = "<makemigrations/migrate/watch>"
