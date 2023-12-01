package models

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/maatko/spotisong/database"
)

type ModelField struct {
	Name string
	Type string
	Properties string
}

type Model struct {
	Table string
	Fields [] ModelField
}

var Models map[string] Model = make(map [string] Model)

func (model Model) Migrate() {
	var builder strings.Builder

	// first the whole table has to be dropped
	builder.WriteString(fmt.Sprintf("DROP TABLE IF EXISTS %v;", model.Table))

	builder.WriteString("CREATE TABLE IF NOT EXISTS ")
	builder.WriteString(model.Table)

	if len(model.Fields) == 0 {
		log.Printf("Migrating the '%v' table...\n", model.Table)
		database.Instance.Exec(builder.String())
		return
	}

	builder.WriteString(" (")
	for index, field := range model.Fields {
		builder.WriteString(fmt.Sprintf("%v %v", field.Name, field.Type))
		if len(field.Properties) > 0 {
			builder.WriteString(" ")
			builder.WriteString(field.Properties)
		}
		if index < len(model.Fields) - 1 {
			builder.WriteString(", ")
		}
	}
	builder.WriteString(")")
	
	log.Printf("Migrating the '%v' table...\n", model.Table)
	database.Instance.Exec(builder.String())
}

func Fetch(definition any) *sql.Rows {
	definitionType := reflect.TypeOf(definition)
	definitionValue := reflect.ValueOf(definition)

	model, ok := Models[definitionType.Name()]
	if !ok {
		log.Fatal("Provided definition does not have a registered model [fetch]")
	}

	numberOfFields := definitionType.NumField()
	if numberOfFields != len(model.Fields) {
		log.Fatal("Provided definition does not match the registered model [field count]")
	}

	var query strings.Builder

	query.WriteString(fmt.Sprintf("SELECT * FROM %v", model.Table))

	appended := false
	for i := 0; i < numberOfFields; i++ {
		fieldType := definitionType.Field(i)
		fieldValue := definitionValue.Field(i)
		modelField := model.Fields[i]

		if strings.ToLower(fieldType.Name) != modelField.Name {
			log.Fatal("Provided definition does not match the registered model [field name]")
		}

		value := fieldValue.Interface()
		if fieldType.Type.Kind() == reflect.String {
			if len(fieldValue.String()) == 0 {
				continue
			}
		} else {
			if fieldValue.Int() == -1 {
				continue
			}
		}

		if appended {
			query.WriteString(fmt.Sprintf(" AND %v='%v'", modelField.Name, value))
		} else {
			query.WriteString(fmt.Sprintf(" WHERE %v='%v'", modelField.Name, value))
		}

		appended = true
	}

	rows, err := database.Instance.Query(query.String())
	if err != nil {
		log.Fatal("Failed to query the database")
	}

	return rows
}

func Insert(definition any) int {
	definitionType := reflect.TypeOf(definition)
	definitionValue := reflect.ValueOf(definition)

	model, ok := Models[definitionType.Name()]
	if !ok {
		log.Fatal("Provided definition does not have a registered model [fetch]")
	}

	numberOfFields := definitionType.NumField()
	if numberOfFields != len(model.Fields) {
		log.Fatal("Provided definition does not match the registered model [field count]")
	}

	var query strings.Builder

	query.WriteString(fmt.Sprintf("INSERT INTO %v (", model.Table))
	for i := 0; i < numberOfFields; i++ {
		fieldType := definitionType.Field(i)
		modelField := model.Fields[i]

		if strings.ToLower(fieldType.Name) != modelField.Name {
			log.Fatal("Provided definition does not match the registered model [field name]")
		}

		if strings.Contains(strings.ToUpper(modelField.Properties), "AUTOINCREMENT") {
			continue
		}

		query.WriteString(modelField.Name)
		if i < numberOfFields - 1 {
			query.WriteString(", ")
		}
	}
	query.WriteString(") VALUES (")
	for i := 0; i < numberOfFields; i++ {
		modelField := model.Fields[i]
		fieldValue := definitionValue.Field(i)

		if strings.Contains(strings.ToUpper(modelField.Properties), "AUTOINCREMENT") {
			continue
		}

		query.WriteString(fmt.Sprintf("'%v'", fieldValue.Interface()))

		if i < numberOfFields - 1 {
			query.WriteString(", ")
		}
	}
	query.WriteString(")")

	result, err := database.Instance.Exec(query.String())
	if err != nil {
		log.Fatal("Failed to insert into the database")
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Fatal("Failed to insert into the database")
	}

	return int(id)
}

func Register(definition any) {
	modelType := reflect.TypeOf(definition)
	modelValue := reflect.ValueOf(definition)

	var fields [] ModelField
	for i := 0; i < modelType.NumField(); i++ {
		fieldType := modelType.Field(i)
		fieldValue := modelValue.Field(i)

		var typeBuilder strings.Builder
		if (fieldType.Type.Kind() == reflect.String) {
			typeBuilder.WriteString(fmt.Sprintf("VARCHAR(%v)", fieldValue.Interface()))
		} else {
			typeBuilder.WriteString("INTEGER")
		}

		field := ModelField {
			Name: strings.ToLower(fieldType.Name),
			Type: typeBuilder.String(),
		}

		if properties, ok := fieldType.Tag.Lookup("properties"); ok {
			field.Properties = properties
		}
		
		fields = append(fields, field)
	}

	modelName := modelType.Name()
	Models[modelName] = Model {
		Table: strings.ToLower(modelName),
		Fields: fields,
	}
}
