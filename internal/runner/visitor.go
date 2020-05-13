package runner

import "github.com/tsouza/yasmim/pkg/command"

type depthFirstAcceptor struct {
	visited map[string]bool
}

func (df *depthFirstAcceptor) Accept(cmd *command.Command, v command.Visitor) (command.VisitorReturnCode, error) {
	df.visited = map[string]bool{}
	return df.accept(cmd, v)
}

func (df *depthFirstAcceptor) accept(cmd *command.Command, v command.Visitor) (command.VisitorReturnCode, error)  {
	if _, exists := df.visited[cmd.Name]; exists {
		return command.VisitorContinue, nil
	}
	df.visited[cmd.Name] = true
	if cmd.OnBefore != nil {
		err := cmd.OnBefore()
		if err != nil {
			return -1, err
		}
	}
	retVal := v.VisitBefore(cmd)
	if retVal >= command.VisitorStop {
		return retVal, nil
	}
	for _, d := range cmd.Dependencies {
		retVal, err := df.accept(d, v)
		if err != nil || retVal == command.VisitorInterrupted {
			return retVal, err
		}
	}
	err := v.VisitAfter(cmd)
	if err != nil {
		return -1, err
	}
	if cmd.OnAfter != nil {
		err := cmd.OnAfter()
		if err != nil {
			return -1, err
		}
	}
	return command.VisitorContinue, nil
}
