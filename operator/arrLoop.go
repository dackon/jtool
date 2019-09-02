package operator

import (
	"fmt"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type loop struct {
}

func (op *loop) do(jv *jvalue.V) (*jvalue.V, error) {
	if jv.JType != jutil.JTArray {
		return nil, fmt.Errorf("operator '$loop': value is not array")
	}

	arrv := jv.Value.([]*jvalue.V)
	if len(arrv) != 2 {
		return nil, fmt.Errorf("operator '$loop': parameter array size must " +
			"be 2. 1st elem must be array, 2nd elem must be operator")
	}

	e1 := arrv[0]
	e2 := arrv[1]
	if e1.JType != jutil.JTArray || e2.JType != jutil.JTString {
		return nil, fmt.Errorf("operator '$loop': parameter array size must " +
			"be 2. 1st elem must be array, 2nd elem must be operator")
	}

	opStr, _ := e2.GetString()
	theOP, ok := gOPMap[opStr]
	if !ok {
		return nil, fmt.Errorf("operator '$loop': '%s' not support", opStr)
	}

	paramArr := e1.Value.([]*jvalue.V)
	valArr := make([]*jvalue.V, 0, len(paramArr))
	for i, v := range paramArr {
		sube, err := theOP.do(v)
		if err != nil {
			return nil, fmt.Errorf("operator '$loop': '%s' failed at %d. "+
				"Err is %s", opStr, i, err)
		}
		valArr = append(valArr, sube)
	}

	newJV := &jvalue.V{
		JType: jutil.JTArray,
		Value: valArr,
	}
	return newJV, nil
}
