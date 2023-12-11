package models

import (
	"spotisong/api"
	"time"
)

type Session struct {
	ID        int       `key:"primary"`
	User      User      `key:"foreign"`
	CreatedAt time.Time `default:"CURRENT_TIMESTAMP"`
	ExpiresAt time.Time
}

func (session *Session) Fetch(keys ...string) error {
	model, err := api.GetModel(*session)
	if err != nil {
		return err
	}

	return model.Fetch(session, keys...)
}

func (session *Session) Save() error {
	model, err := api.GetModel(*session)
	if err != nil {
		return err
	}

	_, err = model.Insert()
	return err
}
