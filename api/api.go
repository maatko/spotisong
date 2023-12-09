package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"text/template"
)

var AppName string
var AppVersion string
var AppDebug bool

// Used for communication with the database,
// migrations are automatically created with the model
var AppModels map[string]Model = map[string]Model{}
var AppMigrations map[string]Migration = map[string]Migration{}

// Stores paths to all directories
// that might need to be queried runtime
var AppDirectories map[string]string = map[string]string{
	"sources":    "./app/",
	"resources":  "./app/resources",
	"migrations": "./app/migrations",
	"templates":  "./app/templates",
}

func InitializeApp(registerModels func() ModelImplementations) error {
	AppName = os.Getenv("APP_NAME")
	AppVersion = os.Getenv("APP_VERSION")

	debug, err := strconv.ParseBool(os.Getenv("APP_DEBUG"))
	if err != nil {
		return fmt.Errorf("`APP_DEBUG` has an invalid value `%v` needs to be a boolean", debug)
	}

	for _, impl := range registerModels() {
		err = RegisterModel(impl)
		if err != nil {
			return err
		}
	}

	AppDebug = debug
	return nil
}

func RenderTemplate(response http.ResponseWriter, data any, statusCode int, paths ...string) error {
	var templates []string
	for _, path := range paths {
		templates = append(templates, GetTemplate(path))
	}

	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		return err
	}

	response.WriteHeader(statusCode)
	err = tmpl.Execute(response, data)
	return err
}

func RegisterModel(impl any) error {
	modelName := reflect.TypeOf(impl).Name()
	if _, has := AppModels[modelName]; has {
		return fmt.Errorf("model '%v' already exists", modelName)
	}

	model := NewModel(impl)
	modelMigration := NewMigration(model)
	if _, ok := AppMigrations[modelName]; ok {
		return fmt.Errorf("'%s' migration already exists", modelName)
	}

	AppModels[modelName] = model
	AppMigrations[modelName] = modelMigration

	return nil
}

func GetModel(impl any) (Model, error) {
	modelName := reflect.TypeOf(impl).Name()
	if model, ok := AppModels[modelName]; ok {
		return model.CreateFields(impl), nil
	}
	return Model{}, fmt.Errorf("model '%v' does not exist", modelName)
}

func GetSource(path string, args ...any) string {
	if sourcesDirectory, ok := AppDirectories["sources"]; ok {
		return (sourcesDirectory + "/" + fmt.Sprintf(path, args...))
	}
	panic(errors.New("sources directory was not specified"))
}

func GetResource(path string, args ...any) string {
	if resourcesDirectory, ok := AppDirectories["resources"]; ok {
		return (resourcesDirectory + "/" + fmt.Sprintf(path, args...))
	}
	panic(errors.New("resources directory was not specified"))
}

func GetMigration(path string, args ...any) string {
	if migrationDirectory, ok := AppDirectories["migrations"]; ok {
		return (migrationDirectory + "/" + fmt.Sprintf(path, args...))
	}
	panic(errors.New("migrations directory was not specified"))
}

func GetTemplate(path string, args ...any) string {
	if templatesDirectory, ok := AppDirectories["templates"]; ok {
		return (templatesDirectory + "/" + fmt.Sprintf(path, args...))
	}
	panic(errors.New("templates directory was not specified"))
}
