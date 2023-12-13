package app

import (
	"encoding/gob"
	"spotisong/api"
	"spotisong/app/models"
	"spotisong/app/routes"
)

func RegisterAppModels() api.ModelImplementations {
	session := models.Session{}
	gob.Register(session)

	return api.ModelImplementations{
		models.User{},
		session,
	}
}

func RegisterAppRoutes() api.HttpRoutes {
	return api.HttpRoutes{
		"/":     routes.Home{},
		"/auth": routes.Authentication,
	}
}
