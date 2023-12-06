package api

import (
	"fmt"
	"reflect"
	"strings"
)

type ModelField struct {
	NativeType reflect.StructField
	NativeValue reflect.Value

	Name string
	Type string

	AutoIncrement bool
	BelongsTo *ModelInformation
}

type ModelInformation struct {
	Name string
	Fields [] ModelField
}

// used when converting from `golang` data types
// to types that the database `sqlite3` recognizes
// TODO :: (match these types more closly)
var ModelTypeConversionPairs map[string] string = map[string] string  {
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

// a cache for all the registered models
// this is used everywhere starting from migrations
// all the way to inserting into and fetching from the database
var ModelRegistry map [string] ModelInformation = map [string] ModelInformation {}

func RegisterModel(definition any) error {
	definitionType := reflect.TypeOf(definition)
	definitionValue := reflect.ValueOf(definition)
	definitionName := definitionType.Name()

	var fields = [] ModelField {}
	for i := 0; i < definitionType.NumField(); i++ {
		fieldType := definitionType.Field(i)
		fieldValue := definitionValue.Field(i)

		field := ModelField {
			NativeType: fieldType,
			NativeValue: fieldValue,
			Name: strings.ToLower(fieldType.Name),
		}

		var typeName string
		if fieldType.Type.Kind() == reflect.Struct {
			typeName = reflect.TypeOf(fieldValue.Interface()).Name()
		} else {
			typeName = fieldType.Type.Kind().String()
		}
		
		tag := fieldType.Tag
		
		if key, ok := tag.Lookup("key"); ok {
			if key == "primary" {
				field.AutoIncrement = true
				field.BelongsTo = nil
			} else if key == "foreign" {
				field.AutoIncrement = false
				if model, ok := ModelRegistry[typeName]; ok {
					field.BelongsTo = &model
				} else {
					return fmt.Errorf("'%v' belongs to a model that was not created '%v'", definitionName, typeName)
				}
			} else {
				return fmt.Errorf("field '%v' has a key property with invalid value '%v'", fieldType.Name, key)
			}
		}
		
		if sqlType, ok := ModelTypeConversionPairs[typeName]; ok {
			field.Type = sqlType
		} else {
			if field.AutoIncrement || field.BelongsTo != nil {
				field.Type = "INTEGER"
			} else {
				return fmt.Errorf("field '%v' in Type '%v' does not have a valid SQL type", fieldType.Name, definitionName)
			}
		}

		fields = append(fields, field)
	}

	ModelRegistry[definitionName] = ModelInformation {
		Name: strings.ToLower(definitionName),
		Fields: fields,
	}

	return nil
}