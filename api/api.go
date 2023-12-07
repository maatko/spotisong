package api

import (
	"fmt"
	"reflect"
)

var Models map [string] Model = map [string] Model {}

func RegisterModel(impl any) error {
	model, modelName := ModelCreate(impl)
	if _, ok := Models[modelName]; ok {
		return fmt.Errorf("'%s' model already exists", modelName)
	}

	Models[modelName] = model
	return nil
}

func GetModel(impl any) (Model, error) {
	implName := reflect.TypeOf(impl).Name()
	if implModel, ok := Models[implName]; ok {
		return implModel.CreateFields(impl), nil
	}
	return Model {}, fmt.Errorf("model with the name of '%v' does not exist", implName)
}