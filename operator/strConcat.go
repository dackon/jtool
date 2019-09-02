package operator

import (
	"fmt"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type concat struct {
}

// do will concat all items to string, e.g.:
// {"s": {"$concat":["abc", "def"]}} -> {"s": "abcdef"}
// {"s": {"$concat":["abc", 1, false, null]}} -> {"s": "abc1falsenull"}
func (op *concat) do(jv *jvalue.V) (*jvalue.V, error) {
	if jv.JType != jutil.JTArray {
		return nil, fmt.Errorf("operator '$concat': value is not array")
	}

	arrv := jv.Value.([]*jvalue.V)
	ret := ""
	str := ""
	for _, v := range arrv {
		switch v.JType {
		case jutil.JTString:
			str, _ = v.GetString()
		case jutil.JTNull:
			str = "null"
		case jutil.JTArray:
			fallthrough
		case jutil.JTObject:
			fallthrough
		case jutil.JTNone:
			return nil, fmt.Errorf("cannot convert %d to string", v.JType)
		default:
			str = fmt.Sprintf("%v", v.Value)
		}
		ret = fmt.Sprintf("%s%s", ret, str)
	}

	newJV := &jvalue.V{JType: jutil.JTString, Value: ret}
	return newJV, nil
}
