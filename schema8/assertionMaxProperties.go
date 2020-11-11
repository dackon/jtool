package schema8

import (
	"context"
	"fmt"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type assertionMaxProperties struct {
	base
	maxProperties int
}

func newAssertionMaxProperties(name string, n *schemaNode) (assertion, error) {
	if n.jtType != jutil.JTInteger {
		return nil, errWithPath(
			fmt.Errorf(
				"value of maxProperties must be a non-negative integer"), n)
	}

	v := n.value.(int64)
	if v < 0 {
		return nil, errWithPath(
			fmt.Errorf(
				"value of maxProperties must be a non-negative integer"), n)
	}

	ast := &assertionMaxProperties{
		base:          base{name, n},
		maxProperties: int(v),
	}
	return ast, nil
}

func (ast *assertionMaxProperties) Valid(ctx context.Context,
	v *jvalue.V) error {

	if v.JType != jutil.JTObject {
		return nil
	}

	m := v.Value.(map[string]*jvalue.V)
	if len(m) <= ast.maxProperties {
		return nil
	}

	return fmt.Errorf("properties count (%d) > maxProperties (%d)",
		len(m), ast.maxProperties)
}
