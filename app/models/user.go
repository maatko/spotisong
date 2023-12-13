package models

import (
	"errors"
	"net/http"
	"spotisong/api"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `key:"primary"`
	Email     string    `max_length:"512"`
	Username  string    `max_length:"255"`
	Password  string    `max_length:"1024"`
	CreatedAt time.Time `default:"CURRENT_TIMESTAMP"`
}

////////////////////////////////////////
// Utility functions
////////////////////////////////////////

func NewUser(email string, username string, password string) *User {
	return &User{
		Email:    email,
		Username: username,
		Password: password,
	}
}

func NewRequestUser(request *http.Request, requireEmail bool, encryptPassword bool) (*User, error) {
	err := request.ParseForm()
	if err != nil {
		return nil, err
	}

	email := request.PostForm.Get("email")
	if len(email) == 0 && requireEmail {
		return nil, errors.New("email not provided")
	}

	username := request.PostForm.Get("username")
	if len(username) == 0 {
		return nil, errors.New("username not provided")
	}

	password := request.PostForm.Get("password")
	if len(password) == 0 {
		return nil, errors.New("password not provided")
	}

	if encryptPassword {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
		if err != nil {
			return nil, err
		}

		password = string(hash)
	}

	return NewUser(email, username, password), nil
}

////////////////////////////////////////
// Database managing functions
////////////////////////////////////////

func (user *User) Load(keys ...string) error {
	return api.FetchModel(user, keys...)
}

func (user *User) Save() error {
	id, err := api.SaveModel(*user)
	if err != nil {
		return err
	}

	user.ID = int(id)
	return nil
}
