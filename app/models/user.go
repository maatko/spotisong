package models

import (
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
