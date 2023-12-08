package api

import (
	"database/sql"
	"fmt"
	"reflect"
)

var DataBase *sql.DB
var Models map[string]Model = map[string]Model{}
var Migrations map[string]Migration = map[string]Migration{}

const MIGRATIONS_DIRECTORY = "./app/migrations"

func RegisterModel(impl any) error {
	model, modelName := ModelCreate(impl)
	if _, ok := Models[modelName]; ok {
		return fmt.Errorf("'%s' model already exists", modelName)
	}

	migration := MigrationCreate(model)
	if _, ok := Migrations[modelName]; ok {
		return fmt.Errorf("'%s' migration already exists", modelName)
	}

	Migrations[modelName] = migration
	Models[modelName] = model
	return nil
}

func GetModel(impl any) (Model, error) {
	implName := reflect.TypeOf(impl).Name()
	if implModel, ok := Models[implName]; ok {
		return implModel.CreateFields(impl), nil
	}
	return Model{}, fmt.Errorf("model with the name of '%v' does not exist", implName)
}
