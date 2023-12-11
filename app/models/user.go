package models

import (
	"errors"
	"net/http"
	"spotisong/api"
	"time"
)

type User struct {
	ID        int       `key:"primary"`
	Email     string    `max_length:"512"`
	Username  string    `max_length:"255"`
	Password  string    `max_length:"1024"`
	CreatedAt time.Time `default:"CURRENT_TIMESTAMP"`
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

func (user *User) Load(keys ...string) error {
	model, err := api.GetModel(*user)
	if err != nil {
		return err
	}

	return model.Fetch(user, keys...)
}

func (user *User) Save() error {
	model, err := api.GetModel(*user)
	if err != nil {
		return err
	}

	id, err := model.Insert()
	if err != nil {
		return err
	}

	user.ID = int(id)
	return nil
}
