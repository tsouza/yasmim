package command

import (
	"context"
	"github.com/tsouza/yasmim/pkg/log"
	"regexp"
)

type Handler func(rt Runtime, log *log.Logger, in, out interface{}) error
type Hook func() error

type Command struct {
	Name 		 string
	Description	 string
	Input   	 interface{}
	Output		 interface{}
	Dependencies []*Command
	OnBefore	 Hook
	OnAfter		 Hook
	Init		 Hook
	Handler		 Handler

	DependencyOf []*regexp.Regexp
}

func (c *Command) String() string {
	return c.Name
}

type Runtime interface {
	Context() 	  context.Context
	Interrupted() bool
}

type Runner func(ctx context.Context, commandName string, in, out interface{}) error

type Filter func(*Command) bool

type Listener interface {
	OnStart(rt Runtime, log *log.Logger, executionSequence []*Command)
	OnBeforeCommand(rt Runtime, log *log.Logger, cmd *Command, values map[string]interface{})
	OnAfterCommand(rt Runtime, log *log.Logger, cmd *Command, values map[string]interface{})
	OnFinish(rt Runtime, log *log.Logger, values map[string]interface{}, err error, interrupted bool)
}

type NoOpListener struct {}
var _ Listener = (*NoOpListener)(nil)
func (l *NoOpListener) OnStart(_ Runtime, _ *log.Logger, _ []*Command) {}
func (l *NoOpListener) OnBeforeCommand(_ Runtime, _ *log.Logger, _ *Command, _ map[string]interface{}) {}
func (l *NoOpListener) OnAfterCommand(_ Runtime, _ *log.Logger, _ *Command, _ map[string]interface{}) {}
func (l *NoOpListener) OnFinish(_ Runtime, _ *log.Logger, _ map[string]interface{}, _ error, _ bool) {}

func NoOpHandler(_ Runtime, _ *log.Logger, _, _ interface{}) error { return nil }

func NoOpFilter(_ *Command) bool { return true }