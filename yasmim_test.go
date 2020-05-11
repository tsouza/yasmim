package yasmim

import (
	"github.com/stretchr/testify/assert"
	"github.com/tsouza/yasmim/pkg/command"
	"github.com/tsouza/yasmim/pkg/log"
	"github.com/tsouza/yasmim/pkg/option"
	"testing"
)

func TestCommands_Simple(t *testing.T) {
	type testInputType  struct { TestInput  string }
	type testOutputType struct { TestOutput string }

	input := testInputType{ "input" }
	output := testOutputType{}

	err := newRunner(
		func(define command.Define) {
			define.Command("test").
				Input(testInputType{}).
				Output(testOutputType{}).
				Handler(func(rt command.Runtime, log *log.Logger, in, out interface{}) error {
					assert.NotPanics(t, func() {
						input := in.(*testInputType)
						output := out.(*testOutputType)

						assert.Equal(t, "input", input.TestInput)

						output.TestOutput = "output"
					})
					return nil
				})
		}).
		Run("test", &input, &output)

	assert.Empty(t, err)
	assert.Equal(t, "output", output.TestOutput)
}

func TestCommands_Filters(t *testing.T) {
	type testType  struct {
		Executed1  bool
		Executed2  bool
		Executed3  bool
	}

	input := testType{}
	output := testType{}

	err := newRunner(
		func(define command.Define) {
			define.Command("test-1").
				Input(testType{}).
				Output(testType{}).
				Dependencies("test-2").
				Handler(func(rt command.Runtime, log *log.Logger, in, out interface{}) error {
					out.(*testType).Executed1 = true
					return nil
				})
		},
		func(define command.Define) {
			define.Command("test-2").
				Input(testType{}).
				Output(testType{}).
				Dependencies("test-3").
				Handler(func(rt command.Runtime, log *log.Logger, in, out interface{}) error {
					out.(*testType).Executed2 = true
					return nil
				})
		},
		func(define command.Define) {
			define.Command("test-3").
				Input(testType{}).
				Output(testType{}).
				Handler(func(rt command.Runtime, log *log.Logger, in, out interface{}) error {
					out.(*testType).Executed3 = true
					return nil
				})
		}).
		With(option.Excludes("test-2")).
		Run("test-1", &input, &output)

	assert.Empty(t, err)
	assert.True(t, output.Executed1)
	assert.False(t, output.Executed2)
	assert.False(t, output.Executed3)
}
