package api

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"reflect"

	"github.com/gorilla/mux"
)

type Message struct {
	Text  string
	Error bool
}

type ProjectInformation struct {
	Name                string
	Directory           string
	StaticDirectory     string
	TemplatesDirectory  string
	MigrationsDirectory string
	DataBase            *sql.DB
	Models              map[string]Model
	Migrations          map[string]Migration
	Router              *mux.Router
	Messages            []Message
}

func ShowMessage(text string, error bool) {
	Project.Messages = append(Project.Messages, Message{
		Text:  text,
		Error: error,
	})
}

func RegisterRoute(path string, route RouteHandler) {
	router := Project.Router
	if !(path == "/" || len(path) == 0) {
		router = mux.NewRouter()
		router.StrictSlash(true)

		Project.Router.Handle(path, router)
	}

	route.SetupRoutes(router)
}

func RenderTemplate(response http.ResponseWriter, data any, statusCode int, paths ...string) error {
	var templates []string

	for _, path := range paths {
		templates = append(templates, Project.Template(path))
	}

	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		return err
	}

	response.WriteHeader(statusCode)

	err = tmpl.Execute(response, data)

	clear(Project.Messages)
	Project.Messages = []Message{}

	return err
}

func RegisterModel(impl any) error {
	model, modelName := ModelCreate(impl)
	if _, ok := Project.Models[modelName]; ok {
		return fmt.Errorf("'%s' model already exists", modelName)
	}

	migration := MigrationCreate(model)
	if _, ok := Project.Migrations[modelName]; ok {
		return fmt.Errorf("'%s' migration already exists", modelName)
	}

	Project.Migrations[modelName] = migration
	Project.Models[modelName] = model
	return nil
}

func GetModel(impl any) (Model, error) {
	implName := reflect.TypeOf(impl).Name()
	if implModel, ok := Project.Models[implName]; ok {
		return implModel.CreateFields(impl), nil
	}
	return Model{}, fmt.Errorf("model with the name of '%v' does not exist", implName)
}

var Project = ProjectInformation{}

func (project *ProjectInformation) Setup(directory string, name string, static string, templates string, migrations string) error {
	project.Name = name
	project.Directory = directory
	project.StaticDirectory = fmt.Sprintf("%s/%s", directory, static)
	project.TemplatesDirectory = fmt.Sprintf("%s/%s", directory, templates)
	project.MigrationsDirectory = fmt.Sprintf("%s/%s", directory, migrations)

	project.Models = map[string]Model{}
	project.Migrations = map[string]Migration{}
	project.Messages = []Message{}

	var err error

	project.DataBase, err = sql.Open("sqlite3", os.Getenv("DATABASE_CONNECTION"))
	if err != nil {
		return err
	}

	project.Router = mux.NewRouter()
	project.Router.StrictSlash(true)

	return nil
}

func (project ProjectInformation) Src(path string, args ...any) string {
	return project.Directory + "/" + fmt.Sprintf(path, args...)
}

func (project ProjectInformation) Static(path string, args ...any) string {
	return project.StaticDirectory + "/" + fmt.Sprintf(path, args...)
}

func (project ProjectInformation) Template(path string, args ...any) string {
	return project.TemplatesDirectory + "/" + fmt.Sprintf(path, args...)
}
