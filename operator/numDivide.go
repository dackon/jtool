package operator

import (
	"fmt"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type divide struct {
}

// do parameter jv must be an array, and the length of it must >= 2. The 1st
// elem in the array will be dividend, the 2nd elem will be divisor, and the
// rest elems in the array are ignored. The result will be a float64.
func (op *divide) do(jv *jvalue.V) (*jvalue.V, error) {
	if jv.JType != jutil.JTArray {
		return nil, fmt.Errorf("operator '$divide': value is not array")
	}

	var v interface{}

	arrv := jv.Value.([]*jvalue.V)

	if len(arrv) < 2 {
		return nil, fmt.Errorf("operator '$divide': array len must be >= 2")
	}

	if (arrv[0].JType != jutil.JTNumber &&
		arrv[0].JType != jutil.JTInteger) ||
		(arrv[1].JType != jutil.JTNumber &&
			arrv[1].JType != jutil.JTInteger) {
		return nil, fmt.Errorf("operator '$divide': bad type")
	}

	// Divisor must not be zero.
	if arrv[1].JType == jutil.JTInteger {
		n, _ := arrv[1].GetInteger()
		if n == 0 {
			return nil, fmt.Errorf("operator '$divide': divisor is 0")
		}
	} else {
		n, _ := arrv[1].GetNumber()
		if jutil.FloatEquals(n, 0.0) {
			return nil, fmt.Errorf("operator '$divide': divisor is 0")
		}
	}

	n1, _ := arrv[0].GetNumber()
	n2, _ := arrv[1].GetNumber()
	v = n1 / n2

	newJV := &jvalue.V{}
	newJV.JType = jutil.JTNumber
	newJV.Value = v

	return newJV, nil
}
