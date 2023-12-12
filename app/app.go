package app

import (
	"spotisong/api"
	"spotisong/app/models"
	"spotisong/app/routes"
)

func RegisterAppModels() api.ModelImplementations {
	return api.ModelImplementations{
		models.User{},
		models.Session{},
	}
}

func RegisterAppRoutes() api.HttpRoutes {
	return api.HttpRoutes{
		"/":     routes.Home{},
		"/auth": routes.Authentication,
	}
}
