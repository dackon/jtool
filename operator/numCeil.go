package operator

import (
	"fmt"
	"math"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type ceil struct {
}

func (op *ceil) do(jv *jvalue.V) (*jvalue.V, error) {
	var v interface{}

	switch jv.JType {
	case jutil.JTInteger:
		v = jv.Value
	case jutil.JTNumber:
		f, _ := jv.GetNumber()
		v = int64(math.Ceil(f))
	default:
		return nil, fmt.Errorf("operator '$ceil': bad value type %d", jv.JType)
	}

	return &jvalue.V{JType: jutil.JTInteger, Value: v}, nil
}
