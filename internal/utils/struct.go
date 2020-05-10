package utils

import (
	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
	"reflect"
)

func FromStructToMap(from interface{}, to map[string]interface{}) {
	structs.FillMap(from, to)
}

func FromMapToStruct(from map[string]interface{}, to interface{}) {
	err := mapstructure.Decode(from, to)
	if err != nil {
		panic(err)
	}
}

func NewStructFrom(s interface{}) interface{} {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return reflect.New(t).Interface()
}
/*
func CopyValues(from interface{}, to interface{}) error {
	fromValue := reflect.ValueOf(from)
	toValue := reflect.ValueOf(to)

	for i := 0; i < toValue.NumField(); i++ {

	}
	return nil
}*/
