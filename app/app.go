package app

import (
	"spotisong/api"
	"spotisong/app/models"

	_ "github.com/mattn/go-sqlite3"
)

func Initialize() {
	api.RegisterModel(models.User{})
	api.RegisterModel(models.Session{})
}
