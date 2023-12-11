package models

import (
	"errors"
	"net/http"
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

func (user *User) FromRequest(request *http.Request) error {
	err := request.ParseForm()
	if err != nil {
		return err
	}

	if !request.Form.Has("username") {
		return errors.New("username was not provided in the request")
	}
	if !request.Form.Has("password") {
		return errors.New("password was not provided in the request")
	}

	user.Username = request.Form.Get("username")
	user.Password = request.Form.Get("password")

	return nil
}
