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

func (session *Session) Save() error {
	model, err := api.GetModel(*session)
	if err != nil {
		return err
	}

	id, err := model.Insert()
	if err != nil {
		return err
	}

	session.ID = int(id)
	return session.Fetch("id")
}

func (session *Session) Fetch(keys ...string) error {
	model, err := api.GetModel(*session)
	if err != nil {
		return err
	}

	rows, err := model.Fetch(*session, keys...)
	if err != nil {
		return err
	}

	// cause we are only fetching a single user
	// we can just access the first one in the list
	rows.Next()

	err = rows.Scan(&session.ID, &session.User.ID, &session.CreatedAt, &session.CreatedAt)
	if err != nil {
		return err
	}

	rows.Close()
	return session.User.Fetch("id")
}
