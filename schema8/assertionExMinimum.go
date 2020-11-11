package schema8

import (
	"context"
	"fmt"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type assertionExMinimum struct {
	base
}

func newAssertionExMinimum(name string, n *schemaNode) (assertion, error) {
	if n.jtType != jutil.JTInteger && n.jtType != jutil.JTNumber {
		return nil, errWithPath(
			fmt.Errorf("value of exclusiveMinimum must be number"), n)
	}

	ast := &assertionExMinimum{base{name, n}}
	return ast, nil
}

func (ast *assertionExMinimum) Valid(ctx context.Context, v *jvalue.V) error {
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
		min := ast.node.value.(int64)
		fmin := float64(min)
		if fv > fmin {
			return nil
		}
		err = fmt.Errorf("value (%fv) <= exclusiveMinimum (%d)", fv, min)
	case jutil.JTNumber:
		min := ast.node.value.(float64)
		if fv > min {
			return nil
		}
		err = fmt.Errorf("value (%f) <= exclusiveMinimum (%f)", fv, min)
	}

	return err
}
