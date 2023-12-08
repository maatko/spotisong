package models

import "time"

type Session struct {
	ID        int       `key:"primary"`
	User      User      `key:"foreign"`
	CreatedAt time.Time `default:"CURRENT_TIMESTAMP"`
	ExpiresAt time.Time
}
