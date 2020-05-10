package runner

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tsouza/yasmim/pkg/command"
	"testing"
)

func TestDepthFirstAcceptor_Accept_Simple(t *testing.T) {
	cmds := command.NewMap(
		func(define command.Define) {
			define.Command("test-1").
				Dependencies("test-2")
		},
		func(define command.Define) {
			define.Command("test-2").
				Dependencies("test-3")
		},
		func(define command.Define) {
			define.Command("test-3")
		},
	)

	testAcceptor := depthFirstAcceptor{}
	testVisitor := &visitListBuilder{}

	test1 := cmds["test-1"]
	retVal, err := testAcceptor.Accept(test1, testVisitor)

	assert.Empty(t, err)
	assert.Equal(t, command.VisitorContinue, retVal)
	assert.Equal(t, []string{
		"before-test-1", "before-test-2", "before-test-3",
		"after-test-3", "after-test-2", "after-test-1",
	}, testVisitor.visitList)
}

func TestDepthFirstAcceptor_Accept_Revisit(t *testing.T) {
	cmds := command.NewMap(
		func(define command.Define) {
			define.Command("test-1").
				Dependencies("test-2", "test-3")
		},
		func(define command.Define) {
			define.Command("test-2").
				Dependencies("test-3")
		},
		func(define command.Define) {
			define.Command("test-3")
		},
	)

	testAcceptor := depthFirstAcceptor{}
	testVisitor := &visitListBuilder{}

	test1 := cmds["test-1"]
	retVal, err := testAcceptor.Accept(test1, testVisitor)

	assert.Empty(t, err)
	assert.Equal(t, command.VisitorContinue, retVal)
	assert.Equal(t, []string{
		"before-test-1", "before-test-2", "before-test-3",
		"after-test-3", "after-test-2", "after-test-1",
	}, testVisitor.visitList)
}

func TestDepthFirstAcceptor_Accept_Cycle(t *testing.T) {
	cmds := command.NewMap(
		func(define command.Define) {
			define.Command("test-1").
				Dependencies("test-2")
		},
		func(define command.Define) {
			define.Command("test-2").
				Dependencies("test-3")
		},
		func(define command.Define) {
			define.Command("test-3").
				Dependencies("test-1")
		},
	)

	testAcceptor := depthFirstAcceptor{}
	testVisitor := &visitListBuilder{}

	test1 := cmds["test-1"]
	retVal, err := testAcceptor.Accept(test1, testVisitor)

	assert.Empty(t, err)
	assert.Equal(t, command.VisitorContinue, retVal)
	assert.Equal(t, []string{
		"before-test-1", "before-test-2", "before-test-3",
		"after-test-3", "after-test-2", "after-test-1",
	}, testVisitor.visitList)

}

type visitListBuilder struct {
	visitList []string
}

func (v *visitListBuilder) VisitBefore(cmd *command.Command) command.VisitorReturnCode {
	v.visitList = append(v.visitList, fmt.Sprintf("before-%v", cmd.Name))
	return command.VisitorContinue
}

func (v *visitListBuilder) VisitAfter(cmd *command.Command) error {
	v.visitList = append(v.visitList, fmt.Sprintf("after-%v", cmd.Name))
	return nil
}