package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"text/template"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

type Message struct {
	Error bool
	Text  string
}

var AppName string
var AppVersion string
var AppDebug bool

// Used for communication with the database,
// migrations are automatically created with the model
var AppModels map[string]Model = map[string]Model{}
var AppMigrations map[string]Migration = map[string]Migration{}

// Used for displaying messages in templates
// through routes by calling api.Message()
var AppMessages []Message = []Message{}

// Used for managing cookies with the between
// the client and the server
var AppCookieStore *sessions.CookieStore

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

	AppCookieStore = sessions.NewCookieStore(
		securecookie.GenerateRandomKey(512), // authentication key
		securecookie.GenerateRandomKey(256), // encryption key
	)

	AppDebug = debug
	return nil
}

func RenderRoute(response http.ResponseWriter, route string, page string, data any) error {
	tmpl, err := template.ParseFiles(
		GetTemplate("base.html"),
		GetTemplate("%s/base.html", route),
		GetTemplate("%s/%s", route, page),
	)

	if err != nil {
		return err
	}

	err = tmpl.Execute(response, map[string]any{
		"messages": AppMessages,
		"data":     data,
	})

	// make sure to clear all the messages
	// that needed to be rendered in the current route
	AppMessages = nil

	return err
}

func MessageError(text string) {
	if AppMessages == nil {
		AppMessages = []Message{}
	}

	AppMessages = append(AppMessages, Message{
		Error: true,
		Text:  text,
	})
}

func MessageInfo(text string) {
	if AppMessages == nil {
		AppMessages = []Message{}
	}

	AppMessages = append(AppMessages, Message{
		Error: false,
		Text:  text,
	})
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
