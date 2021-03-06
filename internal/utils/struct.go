package utils

import (
	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
	"reflect"
)

func FromStructToMap(from interface{}, to map[string]interface{}) {
	fromValue := reflect.ValueOf(from)
	if fromValue.Kind() == reflect.Ptr {
		fromValue = fromValue.Elem()
	}
	for _, name := range structs.Names(from) {
		fValue := fromValue.FieldByName(name)
		to[name] = fValue.Interface()
	}
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
