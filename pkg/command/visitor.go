package command

type Visitor interface {
	VisitBefore(*Command) VisitorReturnCode
	VisitAfter(*Command) error
}

type Acceptor interface {
	Accept(*Command, Visitor) error
}

type VisitorReturnCode int
const (
	VisitorContinue VisitorReturnCode = iota
	VisitorStop
	VisitorInterrupted
)