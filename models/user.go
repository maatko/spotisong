package models

type User struct {
	ID int `properties:"PRIMARY KEY AUTOINCREMENT"`
	Username string
	Password string
}

func CreateUser(username string, password string) User {
	user := User {
		Username: username,
		Password: password,
	}

	// update the id with the inserted 
	// users id
	user.ID = Insert(user)

	return user
}