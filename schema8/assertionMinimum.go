package schema8

import (
	"context"
	"fmt"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type assertionMinimum struct {
	base
}

func newAssertionMinimum(name string, n *schemaNode) (assertion, error) {
	if n.jtType != jutil.JTInteger && n.jtType != jutil.JTNumber {
		return nil, errWithPath(
			fmt.Errorf("value of minimum must be number"), n)
	}

	return &assertionMinimum{base{name, n}}, nil
}

func (ast *assertionMinimum) Valid(ctx context.Context, v *jvalue.V) error {
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
		if fv < fmin {
			err = fmt.Errorf("value (%f) < minimum (%d)", fv, min)
		}

	case jutil.JTNumber:
		fmin := ast.node.value.(float64)
		if fv < fmin {
			err = fmt.Errorf("value (%f) < minimum (%f)", fv, fmin)
		}
	}

	return err
}
