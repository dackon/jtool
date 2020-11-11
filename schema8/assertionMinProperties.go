package schema8

import (
	"context"
	"fmt"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type assertionMinProperties struct {
	base
	minProperties int
}

func newAssertionMinProperties(name string, n *schemaNode) (assertion, error) {
	if n.jtType != jutil.JTInteger {
		return nil, errWithPath(
			fmt.Errorf(
				"value of minProperties must be a non-negative integer"), n)
	}

	v := n.value.(int64)
	if v < 0 {
		return nil, errWithPath(
			fmt.Errorf(
				"value of minProperties must be a non-negative integer"), n)
	}

	ast := &assertionMinProperties{
		base:          base{name, n},
		minProperties: int(v),
	}
	return ast, nil
}

func (ast *assertionMinProperties) Valid(ctx context.Context,
	v *jvalue.V) error {

	if v.JType != jutil.JTObject {
		return nil
	}

	m := v.Value.(map[string]*jvalue.V)
	if len(m) >= ast.minProperties {
		return nil
	}

	return fmt.Errorf("properties count (%d) < minProperties (%d)",
		len(m), ast.minProperties)
}
