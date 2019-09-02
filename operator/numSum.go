package operator

import (
	"fmt"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type sum struct {
}

func (op *sum) do(jv *jvalue.V) (*jvalue.V, error) {
	if jv.JType != jutil.JTArray {
		return nil, fmt.Errorf("operator '$sum': value is not array")
	}

	var v interface{}
	var f float64

	hasFloat := false
	arrv := jv.Value.([]*jvalue.V)

	for _, v := range arrv {
		if !hasFloat && (v.JType == jutil.JTNumber) {
			hasFloat = true
		}

		switch v.JType {
		case jutil.JTInteger:
			n, _ := v.GetInteger()
			f += float64(n)
		case jutil.JTNumber:
			n, _ := v.GetNumber()
			f += n
		default:
			return nil, fmt.Errorf("operator '$sum': value is not int or float")
		}
	}

	newJV := &jvalue.V{}

	if !hasFloat {
		v = int64(f)
		newJV.JType = jutil.JTInteger
	} else {
		v = f
		newJV.JType = jutil.JTNumber
	}

	newJV.Value = v
	return newJV, nil
}
