package schema8

import (
	"context"
	"fmt"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type assertionMaxItems struct {
	base
	maxItems int
}

func newAssertionMaxItems(name string, n *schemaNode) (assertion, error) {
	if n.jtType != jutil.JTInteger {
		return nil, errWithPath(
			fmt.Errorf("value of maxItems must be a non-negative integer"), n)
	}

	v := n.value.(int64)
	if v < 0 {
		return nil, errWithPath(
			fmt.Errorf("value of maxItems must be a non-negative integer"), n)
	}

	ast := &assertionMaxItems{
		base:     base{name, n},
		maxItems: int(v),
	}
	return ast, nil
}

func (ast *assertionMaxItems) Valid(ctx context.Context, v *jvalue.V) error {
	if v.JType != jutil.JTArray {
		return nil
	}

	arr := v.Value.([]*jvalue.V)
	if len(arr) <= ast.maxItems {
		return nil
	}

	return fmt.Errorf("array length (%d) > maxItems (%d)",
		len(arr), ast.maxItems)
}
