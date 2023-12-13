package routes

import (
	"net/http"
	"os"
	"spotisong/api"
	"spotisong/app/models"
	"text/template"
)

func RenderRoute(response http.ResponseWriter, request *http.Request, route string, page string, data any) error {
	baseFile := api.GetTemplate("%s/base.html", route)

	templates := []string{
		api.GetTemplate("base.html"),
		api.GetTemplate("%s/%s", route, page),
	}

	_, err := os.Stat(baseFile)
	if !os.IsNotExist(err) {
		templates = append(templates, api.GetTemplate("%s/base.html", route))
	}

	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		return err
	}

	session, err := models.GetCookieSession(request)

	var user *models.User
	if err == nil {
		user = &session.User
	} else {
		user = nil
	}

	err = tmpl.Execute(response, map[string]any{
		"Messages": api.AppMessages,
		"Data":     data,
		"User":     user,
	})

	// make sure to clear all the messages
	// that needed to be rendered in the current route
	api.AppMessages = nil

	return err
}
