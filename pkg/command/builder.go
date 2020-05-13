package command

import (
	"fmt"
	"sort"
)

type Define interface {
	Command(string) 		Define
	Input(interface{}) 		Define
	Output(interface{}) 	Define
	Dependencies(...string) Define
	SetAsGlobalDependency() Define
	Handler(Handler) 		Define
	OnBefore(Hook)			Define
	OnAfter(Hook)			Define
}

type Builder func(Define)

func NewMap(bs ...Builder) map[string]*Command {
	mDefs := make(map[string]*metadataDef)
	var gDepDefs []string

	for _, b := range bs {
		mDef := &metadataDef{
			md: &Command{},
			deps: []string{},
		}
		b(mDef)
		mDefs[mDef.md.Name] = mDef
		if mDef.md.GlobalDependency {
			gDepDefs = append(gDepDefs, mDef.md.Name)
		}
	}

	if len(gDepDefs) > 0 {
		for _, mDef := range mDefs {
			if !mDef.md.GlobalDependency {
				mDef.deps = append(gDepDefs, mDef.deps...)
			}
		}
	}

	m := make(map[string]*Command)
	for _, mDef := range mDefs {
		order := map[string]int{}
		for idx, depName := range mDef.deps {
			if dep, exists := mDefs[depName]; exists {
				mDef.md.Dependencies = append(mDef.md.Dependencies, dep.md)
				order[depName] = idx
			} else {
				panic(fmt.Errorf("no such dependency %v", depName))
			}
		}
		sort.Slice(mDef.md.Dependencies, func(a, b int) bool {
			return order[mDef.md.Dependencies[a].Name] < order[mDef.md.Dependencies[b].Name]
		})
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

func (m *metadataDef) SetAsGlobalDependency() Define {
	m.md.GlobalDependency = true
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
