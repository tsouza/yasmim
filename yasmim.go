package yasmim

import (
	"context"
	"github.com/tsouza/yasmim/internal/runner"
	"github.com/tsouza/yasmim/pkg/command"
	"github.com/tsouza/yasmim/pkg/log"
	"github.com/tsouza/yasmim/pkg/option"
)

var commands []command.Builder

func Register(builder command.Builder) {
	commands = append(commands, builder)
}

func newRunner(seq ...command.Builder) Runner {
	cmds := command.NewMap(seq...)
	return &builder{&option.Configuration{}, cmds }
}

func New() Runner {
	cmds := commands
	commands = nil
	return newRunner(cmds...)
}

type Runner interface {
	With(opts ...option.Option) Runner
	Run(command string, in, out interface{}) error
}

type builder struct {
	cfg  *option.Configuration
	cmds map[string]*command.Command
}

func (b *builder) With(opts ...option.Option) Runner {
	for _, opt := range opts {
		opt(b.cfg)
	}
	return b
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