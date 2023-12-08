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

func (user *User) Fetch(keys ...string) error {
	model, err := api.GetModel(*user)
	if err != nil {
		return err
	}

	rows, err := model.Fetch(*user, keys...)
	if err != nil {
		return err
	}

	// cause we are only fetching a single user
	// we can just access the first one in the list
	rows.Next()

	defer rows.Close()
	return rows.Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
	)
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
	return user.Fetch("id")
}
