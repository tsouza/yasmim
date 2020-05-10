package utils

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestFromWildcardToRegexp(t *testing.T) {
	var r *regexp.Regexp

	r = fromWildcardToRegexp(t, "test*")
	assert.True(t, r.MatchString("test-1"))
	assert.False(t, r.MatchString("1-test"))

	r = fromWildcardToRegexp(t, "*test")
	assert.False(t, r.MatchString("test-1"))
	assert.True(t, r.MatchString("1-test"))

	r = fromWildcardToRegexp(t, "*test*")
	assert.True(t, r.MatchString("test-1"))
	assert.True(t, r.MatchString("1-test"))

	r = fromWildcardToRegexp(t, "+test+")
	assert.False(t, r.MatchString("test-1"))
	assert.False(t, r.MatchString("1-test"))
}

func fromWildcardToRegexp(t *testing.T, wildcard string) *regexp.Regexp {
	r, err := FromWildcardToRegexp(wildcard)
	if err != nil {
		t.Fatal(err)
	}
	return r
}
