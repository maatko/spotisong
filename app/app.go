package app

import (
	"fmt"
	"spotisong/app/models"
)

func OnRouteRegister() {
	/////////////////////////////////////////
	// Register all your routes here
	/////////////////////////////////////////
}

func OnModelRegister() {
	/////////////////////////////////////////
	// Register all your models here
	/////////////////////////////////////////

	models.User {}.Register()
	models.Post {}.Register()

	//////////////////////////////////////////

	user := models.User {
		Username: "admin",
		Password: "pwd1234",
	}

	err := user.FetchBy("username", "password")
	if err != nil {
		panic(err)
	}

	fmt.Println("=== User ===")
	fmt.Println("> ID:", user.ID)
	fmt.Println("> Username:", user.Username)
	fmt.Println("> Password:", user.Password)
	fmt.Println("> CreatedAt:", user.Created_At)
	fmt.Println("============")
}