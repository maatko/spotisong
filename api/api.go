package api

import "database/sql"

type API struct {
	DataBase *sql.DB
	Address string
	Port int
}

var Instance API = API {}

func (api API) Create(database string) API {
	db, err := sql.Open("sqlite3", database)
	if err != nil {
		panic(err)
	}

	api.DataBase = db
	return api
}