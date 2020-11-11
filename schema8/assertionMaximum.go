package schema8

import (
	"context"
	"fmt"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type assertionMaximum struct {
	base
}

func newAssertionMaximum(name string, n *schemaNode) (assertion, error) {
	if n.jtType != jutil.JTInteger && n.jtType != jutil.JTNumber {
		return nil, errWithPath(
			fmt.Errorf("value of maximum must be number"), n)
	}

	return &assertionMaximum{base{name, n}}, nil
}

func (ast *assertionMaximum) Valid(ctx context.Context, v *jvalue.V) error {
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
		if fv > fmax {
			err = fmt.Errorf("value (%f) > maximum (%d)", fv, max)
		}
	case jutil.JTNumber:
		max := ast.node.value.(float64)
		if fv > max {
			err = fmt.Errorf("value (%f) > maximum (%f)", fv, max)
		}
	}

	return err
}
