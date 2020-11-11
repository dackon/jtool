package schema8

import (
	"context"
	"fmt"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type assertionMinLength struct {
	base
	minLength int
}

func newAssertionMinLength(name string, n *schemaNode) (assertion, error) {
	if n.jtType != jutil.JTInteger {
		return nil, errWithPath(
			fmt.Errorf("value of minLength must be integer"), n)
	}

	ast := &assertionMinLength{
		base:      base{name, n},
		minLength: int(n.value.(int64)),
	}
	return ast, nil
}

func (ast *assertionMinLength) Valid(ctx context.Context, v *jvalue.V) error {
	if v.JType != jutil.JTString {
		return nil
	}

	s, _ := v.GetString()
	if len(s) >= ast.minLength {
		return nil
	}

	return fmt.Errorf("string length (%d) < minLength (%d)",
		len(s), ast.minLength)
}
