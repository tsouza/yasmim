package option

import (
	"context"
	"github.com/tsouza/yasmim/internal/utils"
	"github.com/tsouza/yasmim/pkg/command"
	"github.com/tsouza/yasmim/pkg/log"
	"regexp"
)

type Option func(*Configuration)

func LogDelegate(lg log.DelegateFactory) Option {
	return func(c *Configuration) {
		c.LogDelegate = lg
	}
}

func Context(ctx context.Context) Option {
	return func(c *Configuration) {
		c.Context = ctx
	}
}

func Listener(l command.Listener) Option {
	return func(c *Configuration) {
		c.Listener = l
	}
}

func Excludes(wildcards ...string) Option {
	r := compileWildcards(wildcards)
	return func(c *Configuration) {
		c.Filter = func(cmd *command.Command) bool {
			for _, reg := range r {
				if reg.MatchString(cmd.Name) {
					return false
				}
			}
			return true
		}
	}
}

func compileWildcards(wildcards []string) []*regexp.Regexp {
	var r []*regexp.Regexp
	for _, w := range wildcards {
		reg, err := utils.FromWildcardToRegexp(w)
		if err != nil {
			panic(err)
		}
		r = append(r, reg)
	}
	return r
}

type Configuration struct {
	LogDelegate	log.DelegateFactory
	Context	    context.Context
	Listener	command.Listener
	Filter		command.Filter
}
