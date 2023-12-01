package models

import (
	"fmt"
	"reflect"
)


func Insert(model any) {
	// Get the type of the struct
	personType := reflect.TypeOf(model)

	// Get the value of the struct
	personValue := reflect.ValueOf(model)

	// Iterate through the fields of the struct
	for i := 0; i < personType.NumField(); i++ {
		// Get the field type and value
		fieldType := personType.Field(i)
		fieldValue := personValue.Field(i)

		// Print field name, type, and value
		fmt.Printf("Field Name: %s\n", fieldType.Name)
		fmt.Printf("Field Type: %s\n", fieldType.Type)
		fmt.Printf("Field Value: %v\n", fieldValue.Interface())
		fmt.Println("----------")
	}
}