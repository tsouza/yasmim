package command

import (
	"fmt"
	"github.com/tsouza/yasmim/internal/utils"
	"strings"
)

type Define interface {
	Command(string) 		Define
	Input(interface{}) 		Define
	Output(interface{}) 	Define
	Dependencies(...string) Define
	Handler(Handler) 		Define
	OnBefore(Hook)		Define
	OnAfter(Hook)		Define
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
			if isWildcard(depName) {
				depNameRegexp, err := utils.FromWildcardToRegexp(depName)
				if err != nil {
					panic(err)
				}
				matchedOne := false
				for mDefName, mDef := range mDefs {
					if depNameRegexp.MatchString(mDefName) {
						matchedOne = true
						mDef.md.Dependencies = append(mDef.md.Dependencies, mDef.md)
					}
				}
				if !matchedOne {
					panic(fmt.Errorf("no dependency matched %v", depName))
				}
			} else {
				if dep, exists := mDefs[depName]; exists {
					mDef.md.Dependencies = append(mDef.md.Dependencies, dep.md)
				} else {
					panic(fmt.Errorf("no such dependency %v", depName))
				}
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

func (m *metadataDef) OnBefore(hook Hook) Define {
	m.md.OnBefore = hook
	return m
}

func (m *metadataDef) OnAfter(hook Hook) Define {
	m.md.OnAfter = hook
	return m
}


func isWildcard(name string) bool {
	return strings.Contains(name, "*")
}
