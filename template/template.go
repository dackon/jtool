package template

import (
	"encoding/json"
	"fmt"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

// Template ...
type Template struct {
	tv *tplValue
}

// Debug ...
func Debug(jt *Template) {
	debug("/root", jt.tv)
}

func debug(key string, tv *tplValue) {
	switch tv.valType {
	case jutil.JTObject:
		mv := tv.value.(map[string]*tplValue)
		fmt.Printf("Object key is %s. Value is [%s].\n", key, tv)
		for k, v := range mv {
			debug(fmt.Sprintf("%s/%s", key, k), v)
		}
		return
	case jutil.JTArray:
		arrv := tv.value.([]*tplValue)
		fmt.Printf("Array key is %s. Value is [%s].\n", key, tv)
		for i, v := range arrv {
			debug(fmt.Sprintf("%s/%d", key, i), v)
		}
		return
	}
}

// Parse ...
func Parse(raw json.RawMessage) (*Template, error) {
	jv, err := jvalue.ParseJSON(raw)
	if err != nil {
		return nil, err
	}

	return ParseJValue(jv)
}

// ParseJValue ...
func ParseJValue(jv *jvalue.V) (*Template, error) {
	tv, err := toTPLValue(jv)
	if err != nil {
		return nil, err
	}

	setIterableDim(tv, 0)
	tpl := &Template{tv: tv}
	return tpl, nil
}

func setIterableDim(v *tplValue, dim int) {
	switch v.valType {
	case jutil.JTObject:
		vm := v.value.(map[string]*tplValue)
		for _, subv := range vm {
			setIterableDim(subv, dim)
		}
	case jutil.JTArray:
		arr := v.value.([]*tplValue)
		if len(arr) == 1 {
			if hasIterableWildcardEJP(arr[0], 1) {
				v.isHomoArray = true
				v.dimension = dim + 1
				setIterableDim(arr[0], v.dimension)
			} else {
				v.dimension = 0
				setIterableDim(arr[0], 0)
			}
		} else {
			for _, subv := range arr {
				setIterableDim(subv, dim)
			}
		}
	}
}

func hasIterableWildcardEJP(v *tplValue, needDim int) bool {
	switch v.valType {
	case jutil.JTObject:
		vm := v.value.(map[string]*tplValue)
		for _, subv := range vm {
			if hasIterableWildcardEJP(subv, needDim) {
				return true
			}
		}
	case jutil.JTArray:
		arr := v.value.([]*tplValue)
		if len(arr) == 1 {
			return hasIterableWildcardEJP(arr[0], needDim+1)
		}

		for _, subv := range arr {
			if hasIterableWildcardEJP(subv, needDim) {
				return true
			}
		}
	case jtPointer:
		tp := v.value.(*tplPointer)
		if len(tp.wildcardIdxArr) >= needDim {
			return true
		}
		return false
	default:
		return false
	}

	return false
}

// Execute ...
func (jt *Template) Execute(jpool *JSONPool) (*jvalue.V, error) {
	if jt.tv.valType != jutil.JTObject && jt.tv.valType != jutil.JTArray {
		return nil, fmt.Errorf("bad type %d", jt.tv.valType)
	}

	return execute(jt.tv, jpool)
}

func execute(tv *tplValue, jpool *JSONPool) (*jvalue.V, error) {
	var ii *iterationInfo
	switch tv.valType {
	case jutil.JTObject:
		ii = generateObject(tv, jpool, nil)
		return ii.jv, ii.err
	case jutil.JTArray:
		if tv.isHomoArray {
			ii = generateArraySameItem(tv, jpool, nil)
		} else {
			ii = generateArrayDiffItem(tv, jpool, nil)
		}
		return ii.jv, ii.err
	case jtPointer:
		ii = generateByPointer(tv, jpool, nil)
		return ii.jv, ii.err
	default:
		v := &jvalue.V{
			JType: tv.valType,
			Value: tv.value,
		}
		return v, nil
	}
	return nil, fmt.Errorf("bad type %d", tv.valType)
}

func generateObject(
	tv *tplValue, jpool *JSONPool, idxArr []int) *iterationInfo {

	//	fmt.Printf("generateObject: idxArr is %v. tv is %s\n", idxArr, tv)

	mtv := tv.value.(map[string]*tplValue)
	mjv := make(map[string]*jvalue.V, len(mtv))

	var jv *jvalue.V

	ii := &iterationInfo{}
	var rii *iterationInfo

	for k, v := range mtv {
		//		fmt.Printf("generateObject: key is %s\n", k)

		switch v.valType {
		case jutil.JTObject:
			rii = generateObject(v, jpool, idxArr)
		case jutil.JTArray:
			if v.isHomoArray {
				rii = generateArraySameItem(v, jpool, idxArr)
			} else {
				tarr := v.value.([]*tplValue)
				if len(tarr) == 1 {
					idxArr = nil
				}
				rii = generateArrayDiffItem(v, jpool, idxArr)
			}
		case jtPointer:
			rii = generateByPointer(v, jpool, idxArr)
		default:
			jv = &jvalue.V{
				JType: v.valType,
				Value: v.value,
			}
			rii = nil
		}

		if rii != nil {
			if rii.err == nil {
				mjv[k] = rii.jv
				if ii.dim < rii.dim {
					ii.dim = rii.dim
				}
			}

			if ii.maxArrSize < rii.maxArrSize {
				ii.maxArrSize = rii.maxArrSize
			}
		} else {
			mjv[k] = jv
		}
	}

	if len(mjv) == 0 {
		ii.err = fmt.Errorf("generateObject: no field generatd")
		return ii
	}

	ii.jv = &jvalue.V{
		JType: jutil.JTObject,
		Value: mjv,
	}
	return ii
}

func generateArraySameItem(tv *tplValue, jpool *JSONPool,
	idxArr []int) *iterationInfo {

	//	fmt.Printf("generateArraySameItem: idxArr is %v. tv is %s.\n", idxArr, tv)

	ii := &iterationInfo{}
	itemTpl := tv.value.([]*tplValue)[0]

	var rii *iterationInfo
	var jvArr []*jvalue.V

	if idxArr == nil {
		idxArr = make([]int, 1)
	} else {
		idxArr = append(idxArr, 0)
	}

	last := len(idxArr) - 1

	//	count := 0
	for {
		//		fmt.Printf("generateArraySameItem: Array index is %d\n", count)
		//		count++

		switch itemTpl.valType {
		case jutil.JTArray:
			if itemTpl.isHomoArray {
				rii = generateArraySameItem(itemTpl, jpool, idxArr)
			} else {
				tarr := itemTpl.value.([]*tplValue)
				if len(tarr) == 1 {
					idxArr = nil
				}
				rii = generateArrayDiffItem(itemTpl, jpool, idxArr)
			}
		case jutil.JTObject:
			rii = generateObject(itemTpl, jpool, idxArr)
		case jtPointer:
			rii = generateByPointer(itemTpl, jpool, idxArr)
		default:
			panic(fmt.Sprintf("bad type %d", itemTpl.valType))
		}

		if ii.maxArrSize < rii.maxArrSize {
			ii.maxArrSize = rii.maxArrSize
		}

		if rii.err != nil {
			// If err is not nil, we cannot just break the loop, because some
			// fields in the middle of the array may not exist, e.g.,
			// [{"a":1}, {"b":1}, {"a":2}], if we want get the field 'a' from
			// 2nd of the array, it will be error, but we cannot just break,
			// because the 3rd item has field 'a'. So, we have to check the
			// maxArrSize to break.
			if idxArr[last]+1 >= ii.maxArrSize {
				break
			}
		} else {
			if ii.dim < rii.dim {
				ii.dim = rii.dim
			}

			if rii.dim < tv.dimension {
				break
			}

			jvArr = append(jvArr, rii.jv)
		}
		idxArr[last] = idxArr[last] + 1
	}

	if tv.dimension > 0 && ii.dim < tv.dimension {
		ii.err = fmt.Errorf("generateArraySameItem: maxDim < currDim")
		return ii
	}

	if len(jvArr) == 0 {
		ii.err = fmt.Errorf("generateArraySameItem: no item generated")
		return ii
	}

	ii.jv = &jvalue.V{
		JType: jutil.JTArray,
		Value: jvArr,
	}

	return ii
}

func generateArrayDiffItem(tv *tplValue, jpool *JSONPool,
	idxArr []int) *iterationInfo {

	//	fmt.Printf("generateArrayDiffItem: idxArr is %v. tv is %s.\n", idxArr, tv)

	var jv *jvalue.V
	var jvArr []*jvalue.V
	var rii *iterationInfo

	tplArr := tv.value.([]*tplValue)
	ii := &iterationInfo{}

	for _, v := range tplArr {
		//		fmt.Printf("generateArrayDiffItem: Array index is %d\n", c)
		switch v.valType {
		case jutil.JTArray:
			if v.isHomoArray {
				rii = generateArraySameItem(v, jpool, idxArr)
			} else {
				tarr := v.value.([]*tplValue)
				if len(tarr) == 1 {
					idxArr = nil
				}
				rii = generateArrayDiffItem(v, jpool, idxArr)
			}
		case jutil.JTObject:
			rii = generateObject(v, jpool, idxArr)
		case jtPointer:
			rii = generateByPointer(v, jpool, idxArr)
		default:
			jv = &jvalue.V{
				JType: v.valType,
				Value: v.value,
			}
			rii = nil
		}

		if rii != nil {
			if ii.maxArrSize < rii.maxArrSize {
				ii.maxArrSize = rii.maxArrSize
			}

			if rii.err != nil {
				jvArr = append(jvArr,
					&jvalue.V{JType: jutil.JTNull, Value: nil})
				continue
			}

			if ii.dim < rii.dim {
				ii.dim = rii.dim
			}

			jvArr = append(jvArr, rii.jv)
		} else {
			jvArr = append(jvArr, jv)
		}
	}

	if len(tplArr) == 1 {
		// If tplArr's size is 1, and tv is nonhomogeneous array, the dimension
		// must set to 0 because for the following situation:
		//
		// [{
		//    "a": [{
		//        "a1": 1,
		//        "a2": [{
		//            "a21": {
		//                "a211": "*foo/a/$"
		//            }
		//        }, {
		//            "a22": [
		//                [
		//                    [
		//                        [{
		//                            "a222": "*foo/a/$"
		//                        }]
		//                    ]
		//                ]
		//            ]
		//        }]
		//    }]
		// }]
		// Here, if "*foo/a/$" is valid, field 'a22' always succeeds, so the
		// dimension of 'a22' must set to 0, or we cannot break the loop.
		ii.dim = 0
	}

	ii.jv = &jvalue.V{
		JType: jutil.JTArray,
		Value: jvArr,
	}

	return ii
}

func generateByPointer(tv *tplValue, jpool *JSONPool,
	idxArr []int) *iterationInfo {

	if tv.valType != jtPointer {
		panic("generateByPointer bad logic")
	}

	p := tv.value.(*tplPointer)
	ii := &iterationInfo{
		maxArrSize: getArrayLen(p, jpool, idxArr),
	}

	if len(p.wildcardIdxArr) == 0 {
		ii.jv, ii.err = jpool.GetValue(p.jp)
		return ii
	}

	if len(p.wildcardIdxArr) <= len(idxArr) {
		jp := p.getExJSONPointer()
		for i, v := range p.wildcardIdxArr {
			jp.Keys[v] = fmt.Sprintf("%d", idxArr[i])
		}
		ii.jv, ii.err = jpool.GetValue(jp)
		if ii.err == nil {
			ii.dim = len(p.wildcardIdxArr)
		}

		return ii
	}

	ii.jv, ii.err = getJSONPointerValueArr(p, jpool, idxArr)
	ii.dim = len(p.wildcardIdxArr)
	return ii
}

func getJSONPointerValueArr(tp *tplPointer, jpool *JSONPool,
	idxArr []int) (*jvalue.V, error) {
	if len(idxArr) >= len(tp.wildcardIdxArr) {
		panic("getJSONPointerValueArr bad logic")
	}

	jp := tp.getExJSONPointer()
	for i := 0; i < len(idxArr); i++ {
		jp.Keys[tp.wildcardIdxArr[i]] = fmt.Sprintf("%d", idxArr[i])
	}

	leftIdxArr := make([]int, len(tp.wildcardIdxArr)-len(idxArr))

	var err error
	var jv *jvalue.V
	var values []*jvalue.V

	allIdxArr := make([]int, 0, len(idxArr)+len(leftIdxArr))
	for _, x := range idxArr {
		allIdxArr = append(allIdxArr, x)
	}

L1:
	for {
		leftArrLenArr := make([]int, 0, len(leftIdxArr))
		tmpIdxArr := allIdxArr[0:len(idxArr)]
		for i, v := range leftIdxArr {
			tmpIdxArr = append(tmpIdxArr, v)
			size := getArrayLen(tp, jpool, tmpIdxArr)
			leftArrLenArr = append(leftArrLenArr, size)
			jp.Keys[tp.wildcardIdxArr[len(idxArr)+i]] = fmt.Sprintf("%d", v)
		}

		jv, err = jpool.GetValue(jp)
		if err == nil {
			values = append(values, jv)
		}

		if !getNextIdxArr(leftIdxArr, false) {
			break
		}

		for i, v := range leftIdxArr {
			if v >= leftArrLenArr[i] {
				if !getNextIdxArr(leftIdxArr, true) {
					break L1
				}
			}
		}
	}

	// If err != nil, we still need to iterate array, for example,
	// [{"c":1}, {"d":2}, {"c":3}], it will fail if we try to get 'c'
	// from the 2nd item, but we cannot stop here, we need to continue
	// iterating the array.

	if len(values) == 0 {
		return nil, fmt.Errorf("getJSONPointerValueArr: failed. jp is %s",
			tp.jp)
	}

	jv = &jvalue.V{
		JType: jutil.JTArray,
		Value: values,
	}
	return jv, nil
}

func getNextIdxArr(idxArr []int, overflow bool) bool {
	last := len(idxArr) - 1
	if overflow {
		for i := last; i >= 0; i-- {
			if idxArr[i] != 0 {
				if i == 0 {
					// idxArr no change
					return false
				}

				idxArr[i-1] = idxArr[i-1] + 1
				for j := i; j < len(idxArr); j++ {
					idxArr[j] = 0
				}
				// Index i-1 of idxArr changed.
				return true
			}
		}
		// idxArr no change
		return false
	}

	idxArr[last] = idxArr[last] + 1
	// last index of idxArr changed.
	return true
}

func getArrayLen(p *tplPointer, jpool *JSONPool, idxArr []int) int {
	if len(p.wildcardIdxArr) < len(idxArr) || len(idxArr) == 0 {
		return 0
	}

	jp := p.getExJSONPointer()
	for i := 0; i < len(idxArr)-1; i++ {
		jp.Keys[p.wildcardIdxArr[i]] = fmt.Sprintf("%d", idxArr[i])
	}

	jp.Keys = jp.Keys[0:p.wildcardIdxArr[len(idxArr)-1]]
	return jpool.GetArrayLen(jp)
}
