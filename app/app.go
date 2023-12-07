package app

import (
	"spotisong/api"
	"time"
)

type User struct {
	ID int `key:"primary"`
	Email string `max_length:"512"`
	Username string `max_length:"255"`
	Password string `max_length:"1024"`
	CreatedAt time.Time `default:"TIMESTAMP"`
}

type Post struct {
	ID int `key:"primary"`
	Owner User `key:"foreign"`
	Title string `max_length:"255"`
	Text string `max_length:"1024"`
	CreatedAt time.Time `default:"TIMESTAMP"`
	UpdatedAt time.Time `default:"TIMESTAMP"`
}

func Initialize() {
	api.RegisterModel(User {})
	api.RegisterModel(Post {})
}