package command

import (
	"fmt"
	"sort"
	"strings"
)

type Define interface {
	Command(string) 		Define
	Input(interface{}) 		Define
	Output(interface{}) 	Define
	Dependencies(...interface{}) Define
	Handler(Handler) 		Define
	OnBefore(Hook)			Define
	OnAfter(Hook)			Define
}

type Builder func(Define)
type DependencyMatcher func(string) bool

func NewMap(bs ...Builder) map[string]*Command {
	mDefs := make(map[string]*metadataDef)

	for _, b := range bs {
		mDef := &metadataDef{
			md: &Command{},
			deps: []interface{}{},
		}
		b(mDef)
		mDefs[mDef.md.Name] = mDef
	}

	m := make(map[string]*Command)
	for _, mDef := range mDefs {
		order := map[string]int{}
		for idx, depMatcher := range mDef.deps {
			switch depMatcher.(type) {
			case string:
				depName := depMatcher.(string)
				if dep, exists := mDefs[depName]; exists {
					mDef.md.Dependencies = append(mDef.md.Dependencies, dep.md)
					order[depName] = idx
				} else {
					panic(fmt.Errorf("no such dependency %v", depName))
				}
				break
			default:
				depMatcherFn := DependencyMatcher(depMatcher.(func(string) bool))
				matchedOne := false
				for mmDefName, mmDef := range mDefs {
					if depMatcherFn(mmDefName) {
						matchedOne = true
						mDef.md.Dependencies = append(mDef.md.Dependencies, mmDef.md)
						order[mmDefName] = idx
					}
				}
				if !matchedOne {
					panic(fmt.Errorf("no dependency matched from %v", mDef.md.Name))
				}
				break
			}
		}
		sort.Slice(mDef.md.Dependencies, func(a, b int) bool {
			return order[mDef.md.Dependencies[a].Name] > order[mDef.md.Dependencies[b].Name]
		})
		m[mDef.md.Name] = mDef.md
	}

	return m
}

type metadataDef struct {
	md   *Command
	deps []interface{}
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

func (m *metadataDef) Dependencies(deps ...interface{}) Define {
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
