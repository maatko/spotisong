package models

import "log"

type User struct {
	ID int `properties:"PRIMARY KEY AUTOINCREMENT"`
	Username string
	Password string
}

func (user User) Create() User {
	// update users id with the autogenerated
	// id from the database
	user.ID = Insert(user)
	
	return user
}

func (user User) Fetch() User {
	rows := Fetch(user)
	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Username, &user.Password)
		if err != nil {
			log.Fatal("Failed to fetch user from database")
		}
	}
	return user
}