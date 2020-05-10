package utils

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestFromStructToMap(t *testing.T) {
	testStrVal, testIntVal := "test", 0
	testStruct := &struct{
		Test1 string
		Test2 int
	}{ testStrVal, testIntVal }
	testMap := map[string]interface{}{}

	FromStructToMap(testStruct, testMap)

	assert.NotEmpty(t, testMap)
	assert.Equal(t, map[string]interface{}{
		"Test1": testStrVal,
		"Test2": testIntVal,
	}, testMap)
}

func TestFromMapToStruct(t *testing.T) {
	testStrVal, testIntVal := "test", 0
	testStruct := &struct{
		Test1 string
		Test2 int
	}{}
	testMap := map[string]interface{}{
		"Test1": testStrVal,
		"Test2": testIntVal,
	}

	FromMapToStruct(testMap, testStruct)

	assert.Equal(t, testStrVal, testStruct.Test1)
	assert.Equal(t, testIntVal, testStruct.Test2)
}

func TestNewStructFrom(t *testing.T) {
	type testStructType struct {}
	testStruct := &testStructType{}
	cloneStruct := NewStructFrom(testStruct)

	assert.Equal(t, reflect.TypeOf(testStruct).Elem(), reflect.TypeOf(cloneStruct).Elem())
	assert.NotPanics(t, func() {
		_ = cloneStruct.(*testStructType)
	})
}
