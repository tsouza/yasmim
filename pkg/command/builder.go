package command

import "fmt"

type Define interface {
	Command(string) 		Define
	Input(interface{}) 		Define
	Output(interface{}) 	Define
	Dependencies(...string) Define
	Handler(Handler) 		Define
}

type Builder func(Define)

func NewMap(bs ...Builder) map[string]*Command {
	mDefs := make(map[string]*metadataDef)

	for _, b := range bs {
		mDef := &metadataDef{
			md: &Command{},
			deps: []string{},
		}
		b(mDef)
		mDefs[mDef.md.Name] = mDef
	}

	m := make(map[string]*Command)
	for _, mDef := range mDefs {
		for _, depName := range mDef.deps {
			if dep, exists := mDefs[depName]; exists {
				mDef.md.Dependencies = append(mDef.md.Dependencies, dep.md)
			} else {
				panic(fmt.Errorf("no such dependency %v", depName))
			}
		}
		m[mDef.md.Name] = mDef.md
	}

	return m
}

type metadataDef struct {
	md   *Command
	deps []string
}

func (m *metadataDef) Command(name string) Define {
	m.md.Name = name
	return m
}

func (m *metadataDef) Input(in interface{}) Define {
	m.md.Input = in
	return m
}

func (m *metadataDef) Output(out interface{}) Define {
	m.md.Output = out
	return m
}

func (m *metadataDef) Dependencies(deps ...string) Define {
	m.deps = deps
	return m
}

func (m *metadataDef) Handler(handler Handler) Define {
	m.md.Handler = handler
	return m
}
