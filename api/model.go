package api

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type ModelFieldProperties struct {
	MaxLength int
	Default string
	AutoIncrement bool
	BelongsTo *ModelInformation
}

type ModelField struct {
	NativeType reflect.StructField
	NativeValue reflect.Value

	Name string
	Type string

	Properties ModelFieldProperties
}

type ModelInformation struct {
	Index int
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
var ModelRegistryIndex int = 0

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
			Properties: ModelFieldProperties {},
		}

		var typeName string
		if fieldType.Type.Kind() == reflect.Struct {
			typeName = reflect.TypeOf(fieldValue.Interface()).Name()
		} else {
			typeName = fieldType.Type.Kind().String()
		}
		
		err := field.Properties.Load(fieldType.Tag, typeName, definitionName, fieldType.Name)
		if err != nil {
			return err
		}
		
		if sqlType, ok := ModelTypeConversionPairs[typeName]; ok {
			field.Type = sqlType
		} else {
			if field.Properties.AutoIncrement || field.Properties.BelongsTo != nil {
				field.Type = "INTEGER"
			} else {
				return fmt.Errorf("field '%v' in Type '%v' does not have a valid SQL type", fieldType.Name, definitionName)
			}
		}

		fields = append(fields, field)
	}

	ModelRegistry[definitionName] = ModelInformation {
		Index : ModelRegistryIndex,
		Name: strings.ToLower(definitionName),
		Fields: fields,
	}

	ModelRegistryIndex++
	return nil
}

func (properties *ModelFieldProperties) Load(tag reflect.StructTag, typeName string, definitionName string, fieldName string) error {
	if key, ok := tag.Lookup("key"); ok {
		if key == "primary" {
			properties.AutoIncrement = true
			properties.BelongsTo = nil
		} else if key == "foreign" {
			properties.AutoIncrement = false
			if model, ok := ModelRegistry[typeName]; ok {
				properties.BelongsTo = &model
			} else {
				return fmt.Errorf("'%v' belongs to a model that was not created '%v'", definitionName, typeName)
			}
		} else {
			return fmt.Errorf("field '%v' has a key property with invalid value '%v'", fieldName, key)
		}
	}

	if maxLength, ok := tag.Lookup("max_length"); ok {
		length, err := strconv.Atoi(maxLength)
		if err != nil {
			return fmt.Errorf("field '%v' has to be an valid number value", fieldName)
		}

		properties.MaxLength = length
	}

	if defaultValue, ok := tag.Lookup("default"); ok {
		properties.Default = defaultValue
	}

	return nil
}