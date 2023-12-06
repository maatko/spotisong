package models

import (
	"spotisong/api"
	"time"
)

type User struct {
	// relationships
	ID int `key:"primary"`
	
	// data
	Username string `max_length:"255"`
	Password string `max_length:"1024"`

	// date
	CreatedAt time.Time `default:"CURRENT_TIMESTAMP"`
}

func (user User) Register() {
	err := api.RegisterModel(user)
	if err != nil {
		panic(err)
	}
}