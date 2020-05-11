package utils

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestFromStructToMap(t *testing.T) {
	testStrVal, testIntVal, testNestedVal := "test", 0, &nestedStruct { Test: "test" }
	testStruct := &struct{
		Test1 string
		Test2 int
		Test3 *nestedStruct
	}{ testStrVal, testIntVal, testNestedVal }
	testMap := map[string]interface{}{}

	FromStructToMap(testStruct, testMap)

	assert.NotEmpty(t, testMap)
	assert.Equal(t, map[string]interface{}{
		"Test1": testStrVal,
		"Test2": testIntVal,
		"Test3": testNestedVal,
	}, testMap)
}

func TestFromMapToStruct(t *testing.T) {
	testStrVal, testIntVal, testNestedVal := "test", 0, &nestedStruct { Test: "test" }
	testStruct := &struct{
		Test1 string
		Test2 int
		Test3 *nestedStruct
	}{}
	testMap := map[string]interface{}{
		"Test1": testStrVal,
		"Test2": testIntVal,
		"Test3": testNestedVal,
	}

	FromMapToStruct(testMap, testStruct)

	assert.Equal(t, testStrVal, testStruct.Test1)
	assert.Equal(t, testIntVal, testStruct.Test2)
	assert.Equal(t, testNestedVal, testStruct.Test3)
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

type nestedStruct struct {
	Test string
}
