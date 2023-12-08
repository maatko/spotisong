package app

import (
	"spotisong/api"
	"spotisong/app/models"
	"spotisong/app/routes"

	_ "github.com/mattn/go-sqlite3"
)

func Initialize() {
	api.RegisterModel(models.User{})
	api.RegisterModel(models.Session{})

	///////////////////////////////

	api.RegisterRoute("/", routes.Home{})
	api.RegisterRoute("/", routes.Auth{})

	///////////////////////////////

	user := models.User{
		Email:    "admin@spotisong.com",
		Username: "admin",
		Password: "pwd123",
	}

	err := user.Save()
	if err != nil {
		panic(err)
	}

	session := models.Session{User: user}

	err = session.Save()
	if err != nil {
		panic(err)
	}
}
