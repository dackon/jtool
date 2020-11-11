package schema8

import (
	"context"
	"fmt"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type assertionMinContains struct {
	base
	minContains int
}

func newAssertionMinContains(name string, n *schemaNode) (assertion, error) {
	if n.jtType != jutil.JTInteger {
		return nil, errWithPath(
			fmt.Errorf(
				"value of minContains must be a non-negative integer"), n)
	}

	v := n.value.(int64)
	if v < 0 {
		return nil, errWithPath(
			fmt.Errorf(
				"value of minContains must be a non-negative integer"), n)
	}

	ast := &assertionMinContains{
		base:        base{name, n},
		minContains: int(v),
	}
	return ast, nil
}

func (ast *assertionMinContains) Valid(ctx context.Context, v *jvalue.V) error {
	return nil
}
