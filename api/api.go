package api

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

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

	auth_key, err := strconv.Atoi(os.Getenv("APP_AUTHENTICATION_KEY_LEN"))
	if err != nil {
		return errors.New("authentication key length must be a valid number")
	}

	enc_key, err := strconv.Atoi(os.Getenv("APP_ENCRYPTION_KEY_LEN"))
	if err != nil {
		return errors.New("encprytion key length must a valid number")
	}

	_, err = os.Open("./cookiestore.keys")

	var authKey, encKey string
	if os.IsNotExist(err) {
		authKey = string(securecookie.GenerateRandomKey(auth_key))
		encKey = string(securecookie.GenerateRandomKey(enc_key))

		os.WriteFile("cookiestore.keys", []byte(authKey+":;;::;;:"+encKey), 0755)
	}

	bytes, err := os.ReadFile("cookiestore.keys")
	if err != nil {
		return err
	}

	array := strings.Split(string(bytes), ":;;::;;:")

	AppCookieStore = sessions.NewCookieStore(
		[]byte(array[0]),
		[]byte(array[1]),
	)

	maxAge, err := strconv.Atoi(os.Getenv("APP_SESSION_DURATION"))
	if err != nil {
		return err
	}

	AppCookieStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
	}

	AppDebug = debug
	return nil
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

func TimeCurrent() time.Time {
	return time.Now().Local()
}

func TimeFormat(time time.Time) string {
	return time.Format("2006-01-02 15:04:05")
}
