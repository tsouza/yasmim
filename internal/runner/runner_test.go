package runner

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tsouza/yasmim/pkg/command"
	"github.com/tsouza/yasmim/pkg/log"
	"testing"
)

func TestRun_Simple(t *testing.T) {
	type testInputType struct { TestInput string }
	type testOutputType struct { TestOutput string }

	testCmdName := "test"
	testInputVal, testOutputVal := "test-input", "test-output"
	testInput, testOutput := testInputType{}, testOutputType{}
	cmds := command.NewMap(
		func(define command.Define) {
			define.Command(testCmdName).
				Input(testInput).
				Output(testOutput).
				Handler(func(rt command.Runtime, log *log.Logger, in interface{}, out interface{}) error {
					var input *testInputType
					var output *testOutputType
					assert.NotPanics(t, func() {
						input = in.(*testInputType)
						output = out.(*testOutputType)
					})
					assert.Equal(t, testInputVal, input.TestInput)
					output.TestOutput = testOutputVal
					return nil
				})
		})

	testInput.TestInput = testInputVal

	run := newCmds(cmds)
	err := run(context.Background(), testCmdName, &testInput, &testOutput)

	assert.Empty(t, err)
	assert.Equal(t, testOutputVal, testOutput.TestOutput)
}

func TestRun_Dependency(t *testing.T) {
	type testInputCmd1Type struct { TestOutput2 string }
	type testOutputCmd1Type struct { TestOutput1 string }
	testInputCmd1, testOutputCmd1 := testInputCmd1Type{}, testOutputCmd1Type{}
	testHandlerCmd1 := func(rt command.Runtime, log *log.Logger, in interface{}, out interface{}) error {
		input := in.(*testInputCmd1Type)
		output := out.(*testOutputCmd1Type)
		output.TestOutput1 = "cmd-1<=" + input.TestOutput2
		return nil
	}

	type testInputCmd2Type struct { TestInput2 string }
	type testOutputCmd2Type struct { TestOutput2 string }
	testInputCmd2, testOutputCmd2 := testInputCmd2Type{}, testOutputCmd2Type{}
	testHandlerCmd2 := func(rt command.Runtime, log *log.Logger, in interface{}, out interface{}) error {
		input := in.(*testInputCmd2Type)
		output := out.(*testOutputCmd2Type)
		output.TestOutput2 = "cmd-2<=" + input.TestInput2
		return nil
	}

	cmds := command.NewMap(
		func(define command.Define) {
			define.Command("test-1").
				Input(testInputCmd1).
				Output(testOutputCmd1).
				Dependencies("test-2").
				Handler(testHandlerCmd1)
		},
		func(define command.Define) {
			define.Command("test-2").
				Input(testInputCmd2).
				Output(testOutputCmd2).
				Handler(testHandlerCmd2)
		})

	testInputCmd2.TestInput2 = "outside"

	run := newCmds(cmds)
	err := run(context.Background(), "test-1", &testInputCmd2, &testOutputCmd1)

	assert.Empty(t, err)
	assert.Equal(t, "cmd-1<=cmd-2<=outside", testOutputCmd1.TestOutput1)
}

func TestRun_Listener(t *testing.T) {
	type testType struct  { }

	cmds := command.NewMap(
		func(define command.Define) {
			define.Command("test-1").
				Input(testType{}).
				Output(testType{}).
				Dependencies("test-2").
				Handler(command.NoOpHandler)
		},
		func(define command.Define) {
			define.Command("test-2").
				Input(testType{}).
				Output(testType{}).
				Handler(command.NoOpHandler)
		})

	listenerMethodCallCount := map[string]int{}
	var commandCallOrder []string
	l := &testListener{
		assertOnStart: func(executionSequence []*command.Command) {
			assert.Equal(t, fmt.Sprintf("%v", map[string]string{
				"test-2": "test-2", "test-1": "test-1" }),
				fmt.Sprintf("%v", cmds))
			listenerMethodCallCount["start"]++
		},
		assertOnBeforeCommand: func(commandName string) {
			callName := fmt.Sprintf("before-%v", commandName)
			commandCallOrder = append(commandCallOrder, callName)
			listenerMethodCallCount[callName]++
		},
		assertOnAfterCommand: func(commandName string) {
			callName := fmt.Sprintf("after-%v", commandName)
			commandCallOrder = append(commandCallOrder, callName)
			listenerMethodCallCount[callName]++
		},
		assertOnFinish: func() {
			listenerMethodCallCount["finish"]++
		},
	}
	run := newCmdsWithListener(cmds, l)
	err := run(context.Background(), "test-1", &testType{}, &testType{})

	assert.Empty(t, err)
	assert.Equal(t, map[string]int{
		"start": 1,
		"before-test-1": 1, "before-test-2": 1,
		"after-test-2": 1, "after-test-1": 1,
		"finish": 1,
	}, listenerMethodCallCount)
}

func newCmdsWithListener(cmds map[string]*command.Command, l command.Listener) command.Runner {
	return New(cmds, log.NoOpDelegateFactory, l, command.NoOpFilter)
}

func newCmds(cmds map[string]*command.Command) command.Runner {
	return New(cmds, log.NoOpDelegateFactory, &command.NoOpListener{}, command.NoOpFilter)
}

type testListener struct {
	assertOnStart 		   func([]*command.Command)
	assertOnBeforeCommand  func(string)
	assertOnAfterCommand   func(string)
	assertOnFinish		   func()
}

func (l *testListener) OnStart(rt command.Runtime, log *log.Logger, executionSequence []*command.Command) {
	l.assertOnStart(executionSequence)
}
func (l *testListener) OnBeforeCommand(rt command.Runtime, log *log.Logger, commandName string, values map[string]interface{}) {
	l.assertOnBeforeCommand(commandName)
}
func (l *testListener) OnAfterCommand(rt command.Runtime, log *log.Logger, commandName string, values map[string]interface{}) {
	l.assertOnAfterCommand(commandName)
}
func (l *testListener) OnFinish(rt command.Runtime, log *log.Logger, values map[string]interface{}, err error, interrupted bool) {
	l.assertOnFinish()
}