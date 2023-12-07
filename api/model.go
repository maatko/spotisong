package api

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type ModelFieldProperties struct {
	MaxLength int
	PrimaryKey bool
	Default string
	BelongsTo *Model
}

type ModelField struct {
	Name string
	Type string
	Meta reflect.StructField
	Info reflect.StructTag
	Value reflect.Value
	Properties ModelFieldProperties
}

type Model struct {
	ID int
	Name string
	Fields [] ModelField
}

//////////////////////////
// Type Conversion
// (TODO :: More precise type conversion)
//////////////////////////

var SQL_TYPES map[string] string = map[string] string  {
	"bool": 	"BOOLEAN",
	"string": 	"VARCHAR",
	"uint8": 	"INTEGER",
	"uint16": 	"INTEGER",
	"uint32": 	"INTEGER",
	"uint64": 	"INTEGER",
	"int8": 	"INTEGER",
	"int16": 	"INTEGER",
	"int32": 	"INTEGER",
	"int64": 	"INTEGER",
	"int": 		"INTEGER",
	"float32": 	"FLOAT",
	"float64": 	"FLOAT",
	"float": 	"FLOAT",
	"Time": 	"TIMESTAMP",
}

//////////////////////////
// Model
//////////////////////////

func ModelCreate(impl any) (Model, string) {
	implName := reflect.TypeOf(impl).Name()
	return Model {
		ID: len(Models),
		Name: strings.ToLower(implName),
	}.CreateFields(impl), implName
}

func (model Model) CreateFields(impl any) Model {
	implType := reflect.TypeOf(impl)
	implValue := reflect.ValueOf(impl)

	// make sure to create new slice of fields
	// for the model so it doesn't get appended
	// to the old slice of fields in the model
	model.Fields = [] ModelField {}

	for i := 0; i < implType.NumField(); i++ {
		fieldType := implType.Field(i)
		fieldValue := implValue.Field(i)

		fieldSQLType := "INTEGER"
		if fieldType.Type.Kind() != reflect.Struct {
			typeName := fieldType.Type.Kind().String()
			if sqlType, ok := SQL_TYPES[typeName]; ok {
				fieldSQLType = sqlType
			} else {
				panic(fmt.Sprintf("Type '%v' in '%v' does not have a SQL type", typeName, implType.Name()))
			}
		}

		model.Fields = append(model.Fields, ModelField {
			Name: fieldType.Name,
			Type: fieldSQLType,
			Meta: fieldType,
			Info: fieldType.Tag,
			Value: fieldValue,
		}.ReadProperties())
	}

	return model
}

func (model Model) GetPrimaryField() *ModelField {
	for _, field := range model.Fields {
		if field.Properties.PrimaryKey {
			return &field
		}
	}
	return nil;
}

//////////////////////////
// ModelField
//////////////////////////

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