package schema8

import (
	"context"
	"fmt"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type assertionExMaximum struct {
	base
}

func newAssertionExMaximum(name string, n *schemaNode) (assertion, error) {
	if n.jtType != jutil.JTInteger && n.jtType != jutil.JTNumber {
		return nil, errWithPath(
			fmt.Errorf("value of exclusiveMaximum must be number"), n)
	}

	ast := &assertionExMaximum{base{name, n}}
	return ast, nil
}

func (ast *assertionExMaximum) Valid(ctx context.Context, v *jvalue.V) error {
	fv := 0.0
	switch v.JType {
	case jutil.JTInteger:
		fv = float64(v.Value.(int64))
	case jutil.JTNumber:
		fv = v.Value.(float64)
	default:
		return nil
	}

	var err error
	switch ast.node.jtType {
	case jutil.JTInteger:
		max := ast.node.value.(int64)
		fmax := float64(max)
		if fv < fmax {
			return nil
		}
		err = fmt.Errorf("value (%f) >= exclusiveMaximum (%d)", fv, max)
	case jutil.JTNumber:
		max := ast.node.value.(float64)
		if fv < max {
			return nil
		}
		err = fmt.Errorf("value (%f) >= exclusiveMaximum (%f)", fv, max)
	}

	return err
}
