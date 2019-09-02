package template

import (
	"fmt"

	"github.com/dackon/jtool/ejp"
	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type tplValue struct {
	valType     jutil.JType
	value       interface{}
	isHomoArray bool // is homogeneous array
	dimension   int
}

func (tv *tplValue) String() string {
	return fmt.Sprintf("HArray %v. Dimension %d.", tv.isHomoArray, tv.dimension)
}

func toTPLValue(node *jvalue.V) (*tplValue, error) {
	var newtv *tplValue
	var err error

	tv := &tplValue{
		valType: node.JType,
	}

	switch node.JType {
	case jutil.JTObject:
		mjv := node.Value.(map[string]*jvalue.V)
		mtv := make(map[string]*tplValue, len(mjv))
		for k, v := range mjv {
			newtv, err = toTPLValue(v)
			if err != nil {
				return nil, err
			}
			mtv[k] = newtv
		}
		tv.value = mtv
		return tv, nil
	case jutil.JTArray:
		arrv := node.Value.([]*jvalue.V)
		newarr := make([]*tplValue, 0, len(arrv))
		for _, v := range arrv {
			newtv, err = toTPLValue(v)
			if err != nil {
				return nil, err
			}
			newarr = append(newarr, newtv)
		}
		tv.value = newarr
		return tv, nil
	case jutil.JTString:
		val, _ := node.GetString()
		if len(val) >= 3 && val[0] == '\\' && val[1] == '\\' && val[2] == '*' {
			tv.value = val[2:]
			return tv, nil
		}

		if len(val) == 0 || (len(val) > 0 && val[0] != '*') {
			tv.value = node.Value
			return tv, nil
		}

		// Here, the val MUST be "*..."
		jp, err := ejp.ParseExJSONPointer(val[1:])
		if err != nil {
			return nil, err
		}

		tp := newTPLPointer(jp)
		tv.valType = jtPointer
		tv.value = tp
		return tv, nil
	}

	tv.value = node.Value
	return tv, nil
}
