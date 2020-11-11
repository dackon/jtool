package schema8

import (
	"context"
	"fmt"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type assertionMaxContains struct {
	base
	maxContains int
}

func newAssertionMaxContains(name string, n *schemaNode) (assertion, error) {
	if n.jtType != jutil.JTInteger {
		return nil, errWithPath(
			fmt.Errorf(
				"value of maxContains must be a non-negative integer"), n)
	}

	v := n.value.(int64)
	if v < 0 {
		return nil, errWithPath(
			fmt.Errorf(
				"value of maxContains must be a non-negative integer"), n)
	}

	ast := &assertionMaxContains{
		base:        base{name, n},
		maxContains: int(v),
	}
	return ast, nil
}

func (ast *assertionMaxContains) Valid(ctx context.Context, v *jvalue.V) error {
	return nil
}
