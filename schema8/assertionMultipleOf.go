package schema8

import (
	"context"
	"fmt"
	"math"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type assertionMultipleOf struct {
	base
	fnum float64
}

func newAssertionMultipleOf(name string, n *schemaNode) (assertion, error) {
	fnum := 0.0
	switch n.jtType {
	case jutil.JTInteger:
		v := n.value.(int64)
		if v <= 0 {
			return nil, fmt.Errorf("value of multipleOf must > 0")
		}
		fnum = float64(v)
	case jutil.JTNumber:
		v := n.value.(float64)
		if v <= 0 {
			return nil, fmt.Errorf("value of multipleOf must > 0")
		}
		fnum = v
	default:
		return nil, errWithPath(
			fmt.Errorf("value of multipleOf must be number"), n)
	}

	ast := &assertionMultipleOf{
		base: base{name, n},
		fnum: fnum,
	}
	return ast, nil
}

func (ast *assertionMultipleOf) Valid(ctx context.Context, v *jvalue.V) error {
	if v.JType != jutil.JTInteger && v.JType != jutil.JTNumber {
		return nil
	}

	num := 0.0

	switch v.JType {
	case jutil.JTInteger:
		num = float64(v.Value.(int64))
	case jutil.JTNumber:
		num = v.Value.(float64)
	default:
	}

	if (math.Abs(ast.fnum) > math.Abs(num)) && (num != 0) {
		return fmt.Errorf("value (%f) is not mulitple of (%f)", num, ast.fnum)
	}

	r := num / ast.fnum
	_, frag := math.Modf(r)
	if math.Abs(frag) < jutil.EPSILON {
		return nil
	}

	return fmt.Errorf("value (%f) is not mulitple of (%f)", num, ast.fnum)
}
