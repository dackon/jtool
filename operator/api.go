package operator

import (
	"fmt"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

// RegisterOP ...
func RegisterOP(name string, op Operator) error {
	if len(name) == 0 || len(name) == 1 || name[0] != '$' {
		return fmt.Errorf(
			"operator name must start with '$' and length must > 1")
	}

	_, ok := gOPMap[name]
	if ok {
		return fmt.Errorf("operator %s already exists", name)
	}

	gOPMap[name] = op
	return nil
}

// Do ...
func Do(jv *jvalue.V) (*jvalue.V, error) {
	if jv.JType != jutil.JTObject && jv.JType != jutil.JTArray {
		return nil, fmt.Errorf("operator param must be object or array")
	}

	return do(jv)
}

// do will find the operation key and to call the corresponding function.
// E.g., {"count": {"$add": [1, 2]}} -> {"count": 3}
func do(jv *jvalue.V) (*jvalue.V, error) {
	var err error
	var newJV *jvalue.V

	switch jv.JType {
	case jutil.JTArray:
		arrJV := jv.Value.([]*jvalue.V)
		newArrJV := make([]*jvalue.V, 0, len(arrJV))

		for _, v := range arrJV {
			newJV, err = do(v)
			if err != nil {
				return nil, err
			}
			newArrJV = append(newArrJV, newJV)
		}

		return &jvalue.V{
			JType: jutil.JTArray,
			Value: newArrJV,
		}, nil

	case jutil.JTObject:
		mjv := jv.Value.(map[string]*jvalue.V)
		newMJV := make(map[string]*jvalue.V, len(mjv))

		for k, v := range mjv {
			op, ok := gOPMap[k]
			if ok {
				if len(mjv) > 1 {
					return nil, fmt.Errorf("operation %s has sibling keys", k)
				}

				newJV, err = op.do(v)
				if err != nil {
					return nil, err
				}
				return newJV, nil
			}

			newJV, err = do(v)
			if err != nil {
				return nil, err
			}
			newMJV[k] = newJV
		}

		return &jvalue.V{
			JType: jutil.JTObject,
			Value: newMJV,
		}, nil
	}

	return &jvalue.V{JType: jv.JType, Value: jv.Value}, nil
}
