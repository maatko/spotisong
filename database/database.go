package database

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"
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

func Select(table string, what string, matcher any) (*sql.Rows, error) {
	var query strings.Builder
	query.WriteString(fmt.Sprintf("SELECT %v FROM %v", what, table))

	matcherType := reflect.TypeOf(matcher)
	matcherValue := reflect.ValueOf(matcher)

	numFields := matcherType.NumField()
	if numFields > 0 {
		query.WriteString(" WHERE ")
	}

	var values [] any
	for i := 0; i < numFields; i++ {
		fieldType := matcherType.Field(i)
		fieldValue := matcherValue.Field(i)

		var value any
		if (fieldValue.CanInt()) {
			if (fieldValue.Int() < 0) {
				continue
			}
			value = fieldValue.Int()
		} else {
			if (len(fieldValue.String()) == 0) {
				continue
			}
			value = fieldValue.String()
		}

		if i > 0 {
			query.WriteString(" AND ")	
		}

		query.WriteString(fmt.Sprintf("%v = ?", strings.ToLower(fieldType.Name)))

		values = append(values, value)
	}

	stmt, err := Instance.Prepare(query.String())
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.Query(values...)
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