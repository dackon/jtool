package operator

import (
	"fmt"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type abs struct {
}

func (op *abs) do(jv *jvalue.V) (*jvalue.V, error) {
	var v interface{}

	switch jv.JType {
	case jutil.JTInteger:
		n, _ := jv.GetInteger()
		if n < 0 {
			n = -n
		}
		v = n
	case jutil.JTNumber:
		f, _ := jv.GetNumber()
		if f < 0 {
			f = -f
		}
		v = f
	default:
		return nil, fmt.Errorf("operator '$abs': bad value type %d", jv.JType)
	}

	return &jvalue.V{JType: jv.JType, Value: v}, nil
}
