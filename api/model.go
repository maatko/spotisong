package api

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
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

type ModelImplementations []any

//////////////////////////
// Model
//////////////////////////

func NewModel(impl any) Model {
	return Model{
		ID:   len(AppModels),
		Name: strings.ToLower(reflect.TypeOf(impl).Name()),
	}.CreateFields(impl)
}

func RegisterModel(impl any) error {
	modelName := reflect.TypeOf(impl).Name()
	if _, has := AppModels[modelName]; has {
		return fmt.Errorf("model '%v' already exists", modelName)
	}

	model := NewModel(impl)
	modelMigration := NewMigration(model)
	if _, ok := AppMigrations[modelName]; ok {
		return fmt.Errorf("'%s' migration already exists", modelName)
	}

	AppModels[modelName] = model
	AppMigrations[modelName] = modelMigration

	return nil
}

func GetModel(impl any) (Model, error) {
	modelName := reflect.TypeOf(impl).Name()
	if model, ok := AppModels[modelName]; ok {
		return model.CreateFields(impl), nil
	}
	return Model{}, fmt.Errorf("model '%v' does not exist", modelName)
}

func FetchModel(implPtr any, keys ...string) error {
	model, err := GetModel(reflect.ValueOf(implPtr).Elem().Interface())
	if err != nil {
		return err
	}

	return model.Fetch(implPtr, keys...)
}

func SaveModel(impl any) (int64, error) {
	model, err := GetModel(impl)
	if err != nil {
		return 0, err
	}

	return model.Insert()
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
	sql, values := model.GenerateInsertSQL()

	stmt, err := Server.DataBase.Prepare(sql)
	if err != nil {
		return 0, err
	}

	result, err := stmt.Exec(values...)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (model Model) Fetch(implPtr any, keys ...string) error {
	implType := reflect.TypeOf(implPtr)
	if implType.Kind() != reflect.Ptr {
		return errors.New("pointer to the models implementation must be passed")
	}

	implValue := reflect.ValueOf(implPtr).Elem()

	var values []any
	for _, key := range keys {
		for _, field := range model.Fields {
			if field.Name == key {
				if field.Properties.BelongsTo != nil {
					primaryField := field.Properties.BelongsTo.GetPrimaryField()
					values = append(values, primaryField.GetValue())
				} else {
					values = append(values, field.GetValue())
				}
				break
			}
		}
	}

	var ptrs []any
	for i := 0; i < implValue.NumField(); i++ {
		valueField := implValue.Field(i)
		modelField := model.Fields[i]

		belongsTo := modelField.Properties.BelongsTo
		if belongsTo != nil {
			var nothing interface{}
			ptrs = append(ptrs, &nothing)
			continue
		}

		ptrs = append(ptrs, valueField.Addr().Interface())
	}

	stmt, err := Server.DataBase.Prepare(model.GenerateFetchSQL(keys...))
	if err != nil {
		return err
	}

	row := stmt.QueryRow(values...)
	if row.Err() != nil {
		return row.Err()
	}

	return row.Scan(ptrs...)
}

func (model Model) GenerateInsertSQL() (string, []any) {
	var query strings.Builder
	query.WriteString(fmt.Sprintf("INSERT INTO %v (", model.Name))

	var values []any
	fieldLen := len(model.Fields)
	for _, field := range model.Fields {
		if field.Properties.PrimaryKey || len(field.Properties.Default) > 0 {
			fieldLen--
			continue
		}
		values = append(values, field.GetValue())
	}

	for idx, field := range model.Fields {
		if field.Properties.PrimaryKey || len(field.Properties.Default) > 0 {
			continue
		}

		query.WriteString(field.Name)

		if idx < fieldLen {
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
	return query.String(), values
}

func (model Model) GenerateFetchSQL(keys ...string) string {
	var query strings.Builder
	query.WriteString(fmt.Sprintf("SELECT * FROM %v WHERE ", model.Name))

	for idx, key := range keys {
		query.WriteString(fmt.Sprintf("%v = ?", key))

		if idx < len(keys)-1 {
			query.WriteString(" AND ")
		}
	}

	return query.String()
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
			if field.Meta.Type.Kind() == reflect.Struct {
				if field.Type == "DATETIME" {
					return TimeFormat(field.Value.Interface().(time.Time))
				}
				return field.Value.Interface()
			}
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
