package database

import (
	"database/sql"
	"fmt"
	"log"
)

type TableRowInformation struct {
	ID int
	Name string
	Type string
	NonNull int
	DefaultValue sql.NullString
	PrimaryKey int
}

var Instance *sql.DB

func Initialize(connection string) {
	instance, err := sql.Open("sqlite3", connection)
	if err != nil {
		log.Fatalf("Failed to initialize DataBase '%v'\n", connection)
	}

	Instance = instance
}

func GetTableInformation(table string) ([] TableRowInformation, error)  {
	var information [] TableRowInformation = [] TableRowInformation {}

	stmt, err := Instance.Prepare(fmt.Sprintf("PRAGMA table_info(%s)", table))
	if err != nil {
		return information, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return information, err
	}

	for rows.Next() {
		rowInformation := TableRowInformation {}

		err = rows.Scan(
			&rowInformation.ID,
			&rowInformation.Name,
			&rowInformation.Type,
			&rowInformation.NonNull,
			&rowInformation.DefaultValue,
			&rowInformation.PrimaryKey,
	   )

	   information = append(information, rowInformation)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	
	return information, err
}

func Close() {
	err := Instance.Close()
	if err != nil {
		log.Fatal("Failed to close the DataBase")
	}
}