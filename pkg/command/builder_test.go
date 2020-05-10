package command

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew_Simple(t *testing.T) {
	testName := "test"
	testStruct := struct{ field string }{ testName }
	testNew := NewMap(func(define Define) {
		define.
			Command(testName).
			Input(testStruct).
			Output(testStruct).
			Handler(NoOpHandler)
	})

	assert.Len(t, testNew, 1)
	assert.Contains(t, testNew, testName)
	assert.NotEmpty(t, testNew[testName].Name)
	assert.NotEmpty(t, testNew[testName].Input)
	assert.NotEmpty(t, testNew[testName].Output)
	assert.NotEmpty(t, testNew[testName].Handler)
}

func TestNew_Dependency_Exists(t *testing.T) {
	test1Name := "test-1"
	test2Name := "test-2"
	testStruct := struct{ field string }{ "test" }

	testNew := NewMap(
		func(define Define) {
			define.Command(test1Name).
				Input(testStruct).
				Output(testStruct).
				Handler(NoOpHandler).
				Dependencies(test2Name)
		},
		func(define Define) {
			define.Command(test2Name).
				Input(testStruct).
				Output(testStruct).
				Handler(NoOpHandler).
				Dependencies(test1Name)
		},
	)

	assert.Len(t, testNew, 2)
	assert.Contains(t, testNew, test1Name)
	assert.Contains(t, testNew, test2Name)
	assert.Equal(t, test1Name, testNew[test2Name].Dependencies[0].Name)
	assert.Equal(t, test2Name, testNew[test1Name].Dependencies[0].Name)
}

func TestNew_Dependency_NoExists(t *testing.T) {
	test1Name := "test-1"
	testStruct := struct{ field string }{ "test" }

	assert.Panics(t, func() {
		_ = NewMap(
			func(define Define) {
				define.Command(test1Name).
					Input(testStruct).
					Output(testStruct).
					Handler(NoOpHandler).
					Dependencies("unknown-dep")
			},
		)
	})
}

