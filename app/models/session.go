package models

import (
	"time"
)

type Session struct {
	ID        int  `key:"primary"`
	User      User `key:"foreign"`
	CreatedAt time.Time
	ExpiresAt time.Time
}

func NewSession(user User, expiresIn int) Session {
	currentTime := time.Now().Local()
	return Session{
		User:      user,
		CreatedAt: currentTime,
		ExpiresAt: currentTime.Add(time.Second * time.Duration(expiresIn)),
	}
}
