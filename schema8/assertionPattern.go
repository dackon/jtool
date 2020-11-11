package schema8

import (
	"context"
	"fmt"
	"regexp"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type assertionPattern struct {
	base
	reg    *regexp.Regexp
	regStr string
}

func newAssertionPattern(name string, n *schemaNode) (assertion, error) {
	if n.jtType != jutil.JTString {
		return nil, errWithPath(
			fmt.Errorf("value of pattern must be string"), n)
	}

	s := n.value.(string)
	ap := &assertionPattern{
		base:   base{name, n},
		regStr: s,
	}

	reg, err := regexp.Compile(s)
	if err != nil {
		return nil, errWithPath(fmt.Errorf("Compile pattern failed."), n)
	}
	ap.reg = reg
	return ap, nil
}

func (ast *assertionPattern) Valid(ctx context.Context, v *jvalue.V) error {
	if v.JType != jutil.JTString {
		return nil
	}

	s, _ := v.GetString()
	if ast.reg.MatchString(s) {
		return nil
	}

	return fmt.Errorf("string '%s' doesn't match pattern '%s'", s, ast.regStr)
}
