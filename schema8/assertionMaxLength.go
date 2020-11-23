package schema8

import (
	"context"
	"fmt"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type assertionMaxLength struct {
	base
	maxLength int
}

func newAssertionMaxLength(name string, n *schemaNode) (assertion, error) {
	if n.jtType != jutil.JTInteger {
		return nil, errWithPath(
			fmt.Errorf("value of maxLength must be a string"), n)
	}

	ast := &assertionMaxLength{
		base:      base{name, n},
		maxLength: int(n.value.(int64)),
	}
	return ast, nil
}

func (ast *assertionMaxLength) Valid(ctx context.Context, v *jvalue.V) error {
	if v.JType != jutil.JTString {
		return nil
	}

	s, _ := v.GetString()
	r := []rune(s)
	if len(r) <= ast.maxLength {
		return nil
	}

	return fmt.Errorf("string length (%d) > maxLength (%d)",
		len(s), ast.maxLength)
}
