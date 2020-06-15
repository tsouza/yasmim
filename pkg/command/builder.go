package command

import (
	"fmt"
	"github.com/tsouza/yasmim/internal/utils"
	"regexp"
	"sort"
)

type Define interface {
	Command(string) 		Define
	Description(string)		Define
	Input(interface{}) 		Define
	Output(interface{}) 	Define
	Dependencies(...string) Define
	DependencyOf(...string) Define
	Handler(Handler) 		Define
	Init(Hook) 				Define
	OnBefore(Hook)			Define
	OnAfter(Hook)			Define
}

type Builder func(Define)

func NewMap(bs ...Builder) map[string]*Command {
	mDefs := make(map[string]*metadataDef)
	var gDepDefs []*Command

	for _, b := range bs {
		mDef := &metadataDef{
			md: &Command{},
			deps: []string{},
		}
		b(mDef)
		mDefs[mDef.md.Name] = mDef
		if len(mDef.md.DependencyOf) > 0 {
			gDepDefs = append(gDepDefs, mDef.md)
		}
	}

	if len(gDepDefs) > 0 {
		for _, mDef := range mDefs {
			if len(mDef.md.DependencyOf) == 0 {
				for _, gDepDef := range gDepDefs {
					for _, r := range gDepDef.DependencyOf {
						if r.MatchString(mDef.md.Name) {
							mDef.deps = append([]string{ gDepDef.Name }, mDef.deps...)
						}
					}
				}
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

func (m *metadataDef) Description(description string) Define {
	m.md.Description = description
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

func (m *metadataDef) DependencyOf(deps ...string) Define {
	var depsR []*regexp.Regexp
	for _, dep := range deps {
		depR, err := utils.FromWildcardToRegexp(dep)
		if err != nil {
			panic(err)
		}
		depsR = append(depsR, depR)
	}
	m.md.DependencyOf = depsR
	return m
}

func (m *metadataDef) Handler(handler Handler) Define {
	m.md.Handler = handler
	return m
}

func (m *metadataDef) Init(init Hook) Define {
	m.md.Init = init
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
