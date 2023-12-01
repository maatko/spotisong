package models

type User struct {
	ID int `properties:"PRIMARY KEY AUTOINCREMENT"`
	Email string
	Password string
}

func Create(username string, password string) User {
	user := User {
		ID: 0,
		Password: password,
	}

	Insert(user)
	return user
}