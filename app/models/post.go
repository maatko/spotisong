package models

import (
	"spotisong/api"
	"time"
)

type Post struct {
	// relationships
	ID int 				`key:"primary"`
	Owner User 			`key:"foreign"`

	// data
	Title string 		`max_length:"255"`
	Text string 		`max_length:"1024"`

	// dates
	CreatedAt time.Time `default:"CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `default:"CURRENT_TIMESTAMP"`
}

func (post Post) Register() {
	err := api.RegisterModel(post)
	if err != nil {
		panic(err)
	}
}