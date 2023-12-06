package models

import (
	"spotisong/api"
	"time"
)

type Post struct {
	// relationships
	ID int `key:"primary"`
	Owner User `key:"foreign"`

	// data
	Title string `max_length:"255"`
	Text string `max_length:"1024"`

	// dates
	Created_At time.Time `default:"CURRENT_TIMESTAMP"`
	Updated_At time.Time `default:"CURRENT_TIMESTAMP"`
}

func (post Post) Register() {
	err := api.RegisterModel(post)
	if err != nil {
		panic(err)
	}
}

func (post *Post) Insert() error {
	model, err := api.GetModel(*post)
	if err != nil {
		return err
	}
	
	return model.Insert(*post)
}
