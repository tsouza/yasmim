package yasmim

import (
	"context"
	"fmt"
	"github.com/tsouza/yasmim/internal/runner"
	acceptor2 "github.com/tsouza/yasmim/internal/utils/acceptor"
	"github.com/tsouza/yasmim/pkg/command"
	"github.com/tsouza/yasmim/pkg/log"
	"github.com/tsouza/yasmim/pkg/option"
)

var commands 	   []command.Builder

func Register(builder command.Builder) {
	commands = append(commands, builder)
}

func newRunner(seq ...command.Builder) Yasmim {
	cmds := command.NewMap(seq...)
	return &builder{&option.Configuration{}, cmds }
}

func New() Yasmim {
	cmds := commands
	commands = nil
	return newRunner(cmds...)
}

type Yasmim interface {
	Commands() map[string]*command.Command
	With(opts ...option.Option) Yasmim
	Run(command string, in, out interface{}) error
	Accept(visitor command.Visitor, commandName string) error
}

type builder struct {
	cfg  *option.Configuration
	cmds map[string]*command.Command
}

func (b *builder) Commands() map[string]*command.Command {
	return b.cmds
}

func (b *builder) With(opts ...option.Option) Yasmim {
	for _, opt := range opts {
		opt(b.cfg)
	}
	return b
}


func (b *builder) Accept(visitor command.Visitor, commandName string) error {
	if cmd, exists := b.cmds[commandName]; exists {
		acceptor := acceptor2.DepthFirstAcceptor{}
		_, err := acceptor.Accept(cmd, visitor)
		return err
	}
	return fmt.Errorf("no such command %v", commandName)
}


func (b *builder) Run(command string, in, out interface{}) error {
	b.applyDefaults()
	run := runner.New(b.cmds, b.cfg.LogDelegate, b.cfg.Listener, b.cfg.Filter)
	return run(b.cfg.Context, command, in, out)
}

func (b *builder) applyDefaults() {
	if b.cfg.LogDelegate == nil {
		b.cfg.LogDelegate = log.NoOpDelegateFactory
	}
	if b.cfg.Context == nil {
		b.cfg.Context = context.Background()
	}
	if b.cfg.Listener == nil {
		b.cfg.Listener = &command.NoOpListener{}
	}
	if b.cfg.Filter == nil {
		b.cfg.Filter = command.NoOpFilter
	}
}