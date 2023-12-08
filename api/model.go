package api

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type ModelFieldProperties struct {
	MaxLength  int
	PrimaryKey bool
	Default    string
	BelongsTo  *Model
}

type ModelField struct {
	Name       string
	Type       string
	Meta       reflect.StructField
	Info       reflect.StructTag
	Value      reflect.Value
	Properties ModelFieldProperties
}

type Model struct {
	ID     int
	Name   string
	Fields []ModelField
}

//////////////////////////
// Model
//////////////////////////

func ModelCreate(impl any) (Model, string) {
	implName := reflect.TypeOf(impl).Name()
	return Model{
		ID:   len(Project.Models),
		Name: strings.ToLower(implName),
	}.CreateFields(impl), implName
}

func (model Model) CreateFields(impl any) Model {
	implType := reflect.TypeOf(impl)
	implValue := reflect.ValueOf(impl)

	// make sure to create new slice of fields
	// for the model so it doesn't get appended
	// to the old slice of fields in the model
	model.Fields = []ModelField{}

	for i := 0; i < implType.NumField(); i++ {
		fieldType := implType.Field(i)
		fieldValue := implValue.Field(i)

		typeName := fieldType.Type.Kind().String()
		if fieldType.Type.Kind() == reflect.Struct {
			typeName = reflect.TypeOf(fieldValue.Interface()).Name()
		}

		fieldSQLType := "INTEGER"
		if sqlType, ok := SQL_TYPES[typeName]; ok {
			fieldSQLType = sqlType
		}

		model.Fields = append(model.Fields, ModelField{
			Name:  strings.ToLower(fieldType.Name),
			Type:  fieldSQLType,
			Meta:  fieldType,
			Info:  fieldType.Tag,
			Value: fieldValue,
		}.ReadProperties())
	}

	return model
}

func (model Model) Insert() (int64, error) {
	var query strings.Builder
	var values []any

	query.WriteString(fmt.Sprintf("INSERT INTO %v (", model.Name))

	fieldLen := len(model.Fields)
	for idx, field := range model.Fields {
		if field.Properties.PrimaryKey || len(field.Properties.Default) > 0 {
			fieldLen--
			continue
		}

		query.WriteString(field.Name)

		values = append(values, field.GetValue())

		if idx < fieldLen-1 {
			query.WriteString(", ")
		}
	}
	query.WriteString(") VALUES (")

	valuesLen := len(values)
	for i := 0; i < valuesLen; i++ {
		query.WriteString("?")

		if i < valuesLen-1 {
			query.WriteString(", ")
		}
	}
	query.WriteString(")")

	stmt, err := Project.DataBase.Prepare(query.String())
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(values...)
	if err != nil {
		return 0, err
	}

	insertedID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return insertedID, nil
}

func (model Model) Fetch(impl any, keys ...string) (*sql.Rows, error) {
	var query strings.Builder

	query.WriteString(fmt.Sprintf("SELECT * FROM %v WHERE ", model.Name))

	for idx, key := range keys {
		query.WriteString(fmt.Sprintf("%v = ?", key))

		if idx < len(keys)-1 {
			query.WriteString(" AND ")
		}
	}

	var values []any
	for _, key := range keys {
		for _, field := range model.Fields {
			if field.Name == key {
				values = append(values, field.GetValue())
				break
			}
		}
	}

	stmt, err := Project.DataBase.Prepare(query.String())
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(values...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (model Model) GetPrimaryField() *ModelField {
	for _, field := range model.Fields {
		if field.Properties.PrimaryKey {
			return &field
		}
	}
	return nil
}

//////////////////////////
// ModelField
//////////////////////////

func (field ModelField) GetValue() any {
	if field.Properties.BelongsTo != nil {
		primaryField := field.Properties.BelongsTo.GetPrimaryField()
		if primaryField == nil {
			return nil
		}
		return primaryField.GetValue()
	} else {
		if field.Value.CanInt() {
			return field.Value.Int()
		} else if field.Value.CanFloat() {
			return field.Value.Float()
		} else if field.Meta.Type.Kind() == reflect.Bool {
			return field.Value.Bool()
		} else {
			return field.Value.String()
		}
	}
}

func (field ModelField) ReadProperties() ModelField {
	properties := &field.Properties

	var err error
	if value, ok := field.Info.Lookup("max_length"); ok {
		properties.MaxLength, err = strconv.Atoi(value)
		if err != nil {
			panic(fmt.Sprintf("'%v' has a max length attribute that is not a number (%v)", field.Name, value))
		}
	}

	if value, ok := field.Info.Lookup("key"); ok {
		if value == "primary" {
			properties.PrimaryKey = true
		} else if value == "foreign" {
			properties.PrimaryKey = false

			if field.Meta.Type.Kind() != reflect.Struct {
				panic(fmt.Sprintf("'%v' has a foreign key attribute but is not a struct", field.Name))
			}

			ownerModel, err := GetModel(field.Value.Interface())
			if err != nil {
				panic(err)
			}

			properties.BelongsTo = &ownerModel
		} else {
			panic(fmt.Sprintf("'%v' has a invalid key attribute (%v) <primary/foreign>", field.Name, value))
		}
	}

	properties.Default = field.Info.Get("default")
	return field
}
