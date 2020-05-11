package runner

import (
	"context"
	"fmt"
	"github.com/tsouza/yasmim/internal/utils"
	"github.com/tsouza/yasmim/pkg/command"
	"github.com/tsouza/yasmim/pkg/log"
)

func New(cmds map[string]*command.Command, newLogDelegate log.DelegateFactory, listener command.Listener, filter command.Filter) command.Runner {
	return func(ctx context.Context, commandName string, in interface{}, out interface{}) error {
		if cmd, exists := cmds[commandName]; exists {
			rt := &runtime{ ctx }
			acceptor := depthFirstAcceptor{}
			counter := &commandCounter{rt, filter, 0}
			retVal, err := acceptor.Accept(cmd, counter)
			if err != nil || retVal == command.VisitorInterrupted {
				return err
			}
			values := map[string]interface{}{}
			listenerLogger := log.NewLogger(newLogDelegate("listener"))
			listener.OnStart(rt, listenerLogger, counter.totalCommands)
			if in != nil {
				utils.FromStructToMap(in, values)
			}
			executor := &commandExecutor{
				rt:             rt,
				f:				filter,
				listener:       listener,
				listenerLogger: listenerLogger,
				newLogDelegate: newLogDelegate,
				values:         values,
			}
			retVal, err = acceptor.Accept(cmd, executor)
			if err != nil || retVal >= command.VisitorStop {
				listener.OnFinish(rt, listenerLogger, values, err, retVal == command.VisitorInterrupted)
				return err
			}
			if out != nil {
				utils.FromMapToStruct(values, out)
			}
			listener.OnFinish(rt, listenerLogger, values, nil, false)
			return nil
		}
		return fmt.Errorf("no such command %v", commandName)
	}
}

type runtime struct {
	ctx		context.Context
}

func (rt *runtime) Context() context.Context {
	return rt.ctx
}

func (rt *runtime) Interrupted() bool {
	return rt.ctx.Err() != nil
}

type commandExecutor struct {
	rt 		 		command.Runtime
	f				command.Filter
	listener 		command.Listener
	listenerLogger  *log.Logger
	newLogDelegate  log.DelegateFactory
	values			map[string]interface{}
}

func (v *commandExecutor) VisitBefore(cmd *command.Command) command.VisitorReturnCode {
	return doCallVisitor(v.rt, v.f, cmd, func() {
		v.listener.OnBeforeCommand(v.rt, v.listenerLogger, cmd.Name, v.values)
	})
}

func (v *commandExecutor) VisitAfter(cmd *command.Command) error {
	if cmd.Handler != nil {
		logger := log.NewLogger(v.newLogDelegate(cmd.Name))
		in := utils.NewStructFrom(cmd.Input)
		utils.FromMapToStruct(v.values, in)
		out := utils.NewStructFrom(cmd.Output)
		err := cmd.Handler(v.rt, logger, in, out)
		if err != nil {
			return err
		}
		utils.FromStructToMap(out, v.values)
	}
	v.listener.OnAfterCommand(v.rt, v.listenerLogger, cmd.Name, v.values)
	return nil
}

type commandCounter struct {
	rt				command.Runtime
	f				command.Filter
	totalCommands	int
}

func (v *commandCounter) VisitBefore(cmd *command.Command) command.VisitorReturnCode {
	return doCallVisitor(v.rt, v.f, cmd, func() {
		v.totalCommands++
	})
}

func (v *commandCounter) VisitAfter(_ *command.Command) error { return nil }

func doCallVisitor(rt command.Runtime, f command.Filter, cmd *command.Command, v func()) command.VisitorReturnCode {
	if rt.Interrupted() {
		return command.VisitorInterrupted
	}
	if !f(cmd) {
		return command.VisitorStop
	}
	v()
	return command.VisitorContinue
}