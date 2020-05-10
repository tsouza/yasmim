package utils

import (
	"fmt"
	"regexp"
)

func FromWildcardToRegexp(wildcard string) (*regexp.Regexp, error) {
	r := ""
	for _, c := range wildcard {
		if c == '*' {
			r = fmt.Sprintf("%v.*?", r)
		} else if c == '+' {
			r = fmt.Sprintf("%v.+?", r)
		} else {
			r = fmt.Sprintf("%v%v", r, regexp.QuoteMeta(string(c)))
		}
	}
	return regexp.Compile(fmt.Sprintf("^%v$", r))
}
