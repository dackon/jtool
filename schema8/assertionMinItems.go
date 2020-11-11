package schema8

import (
	"context"
	"fmt"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type assertionMinItems struct {
	base
	minItems int
}

func newAssertionMinItems(name string, n *schemaNode) (assertion, error) {
	if n.jtType != jutil.JTInteger {
		return nil, errWithPath(
			fmt.Errorf("value of minItems must be a non-negative integer"), n)
	}

	v := n.value.(int64)
	if v < 0 {
		return nil, errWithPath(
			fmt.Errorf("value of minItems must be a non-negative integer"), n)
	}

	ast := &assertionMinItems{
		base:     base{name, n},
		minItems: int(v),
	}
	return ast, nil
}

func (ast *assertionMinItems) Valid(ctx context.Context, v *jvalue.V) error {
	if v.JType != jutil.JTArray {
		return nil
	}

	arr := v.Value.([]*jvalue.V)
	if len(arr) >= ast.minItems {
		return nil
	}

	return fmt.Errorf("array length (%d) < minItems (%d)",
		len(arr), ast.minItems)
}
