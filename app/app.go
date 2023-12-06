package app

import (
	"spotisong/api"
	"time"
)

type User struct {
	ID int 				`key:"primary"`
	Username string 	`size:"255"`
	Password string 	`size:"1024"`
	CreatedAt time.Time `default:"CURRENT_TIMESTAMP"`
}

type Post struct {
	ID int 				`key:"primary"`
	Owner User 			`key:"foreign"`
	Title string 		`size:"255"`
	Text string 		`size:"1024"`
	CreatedAt time.Time `default:"CURRENT_TIMESTAMP"`
	UpdatedAt time.Time
}

func OnRouteRegister() {
	/////////////////////////////////////////
	// Register all your routes here
	/////////////////////////////////////////
}

func OnModelRegister() {
	/////////////////////////////////////////
	// Register all your models here
	/////////////////////////////////////////

	err := api.RegisterModel(User {})
	if err != nil {
		panic(err)
	}

	err = api.RegisterModel(Post {})
	if err != nil {
		panic(err)
	}
}