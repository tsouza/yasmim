package command

import (
	"context"
	"github.com/tsouza/yasmim/pkg/log"
)

type Handler func(rt Runtime, log *log.Logger, in, out interface{}) error

type Command struct {
	Name 		 string
	Input   	 interface{}
	Output		 interface{}
	Dependencies []*Command
	Handler		 Handler
}

type Runtime interface {
	Context() 	  context.Context
	Interrupted() bool
}

type Runner func(ctx context.Context, commandName string, in, out interface{}) error

type Filter func(*Command) bool

type Listener interface {
	OnStart(rt Runtime, log *log.Logger, totalCommandsToExecute int)
	OnBeforeCommand(rt Runtime, log *log.Logger, commandName string, values map[string]interface{})
	OnAfterCommand(rt Runtime, log *log.Logger,commandName string, values map[string]interface{})
	OnFinish(rt Runtime, log *log.Logger, values map[string]interface{}, err error, interrupted bool)
}

type NoOpListener struct {}
var _ Listener = (*NoOpListener)(nil)
func (l *NoOpListener) OnStart(_ Runtime, _ *log.Logger, _ int) {}
func (l *NoOpListener) OnBeforeCommand(_ Runtime, _ *log.Logger, _ string, _ map[string]interface{}) {}
func (l *NoOpListener) OnAfterCommand(_ Runtime, _ *log.Logger, _ string, _ map[string]interface{}) {}
func (l *NoOpListener) OnFinish(_ Runtime, _ *log.Logger, _ map[string]interface{}, _ error, _ bool) {}

func NoOpHandler(_ Runtime, _ *log.Logger, _, _ interface{}) error { return nil }

func NoOpFilter(_ *Command) bool { return true }