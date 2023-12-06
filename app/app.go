package app

import (
	"fmt"
	"spotisong/app/models"
)

func OnRouteRegister() {
	/////////////////////////////////////////
	// Register all your routes here
	/////////////////////////////////////////

	user := models.User {
		Username: "admin",
		Password: "admin",
	}

	fmt.Println(user.ID)
}

func OnModelRegister() {
	/////////////////////////////////////////
	// Register all your models here
	/////////////////////////////////////////

	models.User {}.Register()
	models.Post {}.Register()
}