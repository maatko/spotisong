package models

import (
	"time"
)

type User struct {
	ID        int       `key:"primary"`
	Email     string    `max_length:"512"`
	Username  string    `max_length:"255"`
	Password  string    `max_length:"1024"`
	CreatedAt time.Time `default:"CURRENT_TIMESTAMP"`
}

func NewUser(username string, password string) User {
	return User{
		Username: username,
		Password: password,
	}
}
