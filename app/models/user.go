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
	Created_At time.Time `default:"CURRENT_TIMESTAMP"`
}

func (user *User) Insert() error {
	model, err := api.GetModel(*user)
	if err != nil {
		return err
	}
	
	return model.Insert(*user)
}

func (user *User) FetchBy(tags ...string)  error {
	model, err := api.GetModel(*user)
	if err != nil {
		return err
	}

	rows, err := model.FetchBy(*user, tags...)
	if err != nil {
		return err
	}

	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Created_At)
		if err != nil {
			return err
		}
	}
	
	return nil
}

func (user User) Register() User {
	err := api.RegisterModel(user)
	if err != nil {
		panic(err)
	}
	return user
}