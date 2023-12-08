package api

import (
	"database/sql"
	"fmt"
	"reflect"
)

type ProjectInformation struct {
	Name                string
	Directory           string
	StaticDirectory     string
	MigrationsDirectory string
}

var DataBase *sql.DB
var Models map[string]Model = map[string]Model{}
var Migrations map[string]Migration = map[string]Migration{}
var Project = ProjectInformation{}

func (project *ProjectInformation) Setup(directory string, name string, static string, migrations string) {
	project.Name = name
	project.Directory = directory
	project.StaticDirectory = fmt.Sprintf("%s/%s", directory, static)
	project.MigrationsDirectory = fmt.Sprintf("%s/%s", directory, migrations)
}

func (project *ProjectInformation) Src(path string, args ...any) string {
	return project.Directory + "/" + fmt.Sprintf(path, args...)
}

func (project *ProjectInformation) Static(path string, args ...any) string {
	return project.StaticDirectory + "/" + fmt.Sprintf(path, args...)
}

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
