package jvalue

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/dackon/jtool/ejp"
	"github.com/dackon/jtool/jutil"
)

// V ...
type V struct {
	JType jutil.JType
	Value interface{}
}

// trimRaw deletes the leading and trailing spaces.
func trimRaw(raw json.RawMessage) json.RawMessage {
	if len(raw) == 0 {
		return raw
	}

	i := 0
	for ; i < len(raw); i++ {
		if jutil.IsSpace(raw[i]) {
			continue
		} else {
			break
		}
	}

	if i > 0 {
		raw = raw[i:]
	}

	if len(raw) == 0 {
		return raw
	}

	i = len(raw) - 1
	for ; i >= 0; i-- {
		if jutil.IsSpace(raw[i]) {
			continue
		} else {
			break
		}
	}

	raw = raw[:i+1]
	return raw
}

// ParseJSON parses json value, e.g., number/string/null/object/array
// to JSONValue.
func ParseJSON(raw json.RawMessage) (*V, error) {
	raw = trimRaw(raw)
	if len(raw) == 0 {
		return nil, ErrBadParam
	}

	jv := &V{}
	switch raw[0] {
	case '{':
		jv.JType = jutil.JTObject
		val := make(map[string]json.RawMessage)
		err := json.Unmarshal(raw, &val)
		if err != nil {
			return nil, err
		}

		mv := make(map[string]*V, len(val))
		for kval, vval := range val {
			tv, err := ParseJSON(vval)
			if err != nil {
				return nil, err
			}

			mv[kval] = tv
		}

		jv.Value = mv
	case '[':
		jv.JType = jutil.JTArray
		val := []json.RawMessage{}
		err := json.Unmarshal(raw, &val)
		if err != nil {
			return nil, err
		}

		av := make([]*V, 0, len(val))
		for _, vval := range val {
			tv, err := ParseJSON(vval)
			if err != nil {
				return nil, err
			}
			av = append(av, tv)
		}
		jv.Value = av
	case '"':
		// Do not need to check whether the last char is '"' or not, because
		// the raw is from json.Unmarshal().
		jv.JType = jutil.JTString
		jv.Value = string(raw[1 : len(raw)-1])
	case 't', 'f':
		jv.JType = jutil.JTBoolean
		str := string(raw)
		if str == "true" {
			jv.Value = true
		} else if str == "false" {
			jv.Value = false
		} else {
			return nil, fmt.Errorf("bad boolean value %s", str)
		}
	case 'n':
		jv.JType = jutil.JTNull
		if string(raw) != "null" {
			return nil, fmt.Errorf("bad null %s", raw)
		}
	default:
		i, err := strconv.ParseInt(string(raw), 10, 64)
		if err == nil {
			jv.JType = jutil.JTInteger
			jv.Value = i
			return jv, nil
		}

		f, err := strconv.ParseFloat(string(raw), 64)
		if err != nil {
			return nil, err
		}
		jv.JType = jutil.JTNumber
		jv.Value = f
	}

	return jv, nil
}

// IsEqual ...
func (v *V) IsEqual(t *V) error {
	if v.JType != t.JType {
		return ErrTypeNotEqual
	}

	switch v.JType {
	case jutil.JTNull:
		return nil
	case jutil.JTBoolean:
		return v.compareBool(t)
	case jutil.JTNumber:
		return v.compareNumber(t)
	case jutil.JTString:
		return v.compareString(t)
	case jutil.JTObject:
		return v.compareObject(t)
	case jutil.JTArray:
		return v.compareArray(t)
	case jutil.JTInteger:
		return v.compareInteger(t)
	default:
		panic("Unknown value type")
	}
}

func (v *V) compareNull(t *V) error {
	return nil
}

func (v *V) compareBool(t *V) error {
	x, _ := v.Value.(bool)
	y, _ := t.Value.(bool)

	if x != y {
		return ErrBoolNotEqual
	}
	return nil
}

func (v *V) compareNumber(t *V) error {
	x, _ := v.Value.(float64)
	y, _ := t.Value.(float64)

	if !jutil.FloatEquals(x, y) {
		return ErrNumberNotEqual
	}

	return nil
}

func (v *V) compareInteger(t *V) error {
	x, _ := v.Value.(int64)
	y, _ := t.Value.(int64)

	if x != y {
		return ErrIntegerNotEqual
	}
	return nil
}

func (v *V) compareString(t *V) error {
	x, _ := v.Value.(string)
	y, _ := t.Value.(string)

	if x != y {
		return ErrStringNotEqual
	}
	return nil
}

func (v *V) compareObject(t *V) error {
	x, _ := v.Value.(map[string]*V)
	y, _ := t.Value.(map[string]*V)

	if len(x) != len(y) {
		return ErrObjectNotEqual
	}

	var err error
	for k, j := range x {
		m, ok := y[k]
		if !ok {
			return ErrObjectNotEqual
		}

		err = j.IsEqual(m)
		if err != nil {
			return err
		}
	}

	return nil
}

func (v *V) compareArray(t *V) error {
	x, _ := v.Value.([]*V)
	y, _ := t.Value.([]*V)

	if len(x) != len(y) {
		return ErrArrayNotEqual
	}

	var err error
	for i := 0; i < len(x); i++ {
		err = x[i].IsEqual(y[i])
		if err != nil {
			return err
		}
	}

	return nil
}

// GetNumber ...
func (v *V) GetNumber() (float64, error) {
	switch v.JType {
	case jutil.JTNumber:
		return v.Value.(float64), nil
	case jutil.JTInteger:
		return float64(v.Value.(int64)), nil
	}

	return 0, fmt.Errorf("value type is %d, not number", v.JType)
}

// GetInteger ...
func (v *V) GetInteger() (int64, error) {
	switch v.JType {
	case jutil.JTInteger:
		return v.Value.(int64), nil
	}

	return 0, fmt.Errorf("value type is %d, not integer", v.JType)
}

// GetString ...
func (v *V) GetString() (string, error) {
	switch v.JType {
	case jutil.JTString:
		return v.Value.(string), nil
	}
	return "", fmt.Errorf("value type is %d, not string", v.JType)
}

// GetStringArr ...
func (v *V) GetStringArr() ([]string, error) {
	switch v.JType {
	case jutil.JTArray:
		arr, _ := v.Value.([]*V)
		sarr := make([]string, 0, len(arr))
		for _, v := range arr {
			s, err := v.GetString()
			if err != nil {
				return nil, err
			}
			sarr = append(sarr, s)
		}
		return sarr, nil
	}

	return nil, fmt.Errorf("value type is %d, not array", v.JType)
}

// GetBool ...
func (v *V) GetBool() (bool, error) {
	switch v.JType {
	case jutil.JTBoolean:
		return v.Value.(bool), nil
	}

	return false, fmt.Errorf("value type is %d, not bool", v.JType)
}

// GetSubJSONValue get JSON value by JSON pointer, please note:
// For /a/b/c, if a or b doesn't exist, the JSON pointer is invalid,
// if a & b are valid, c doesn't exist, this function will return ErrNotFound,
// and up layer can use this error type to define default value.
func (v *V) GetSubJSONValue(jp *ejp.ExJSONPointer) (*V, error) {
	// Here, len(jp.Keys) MUST be >= 1
	if len(jp.Keys) == 1 {
		return v, nil
	}

	jv := v
	// The first key is "", which represents the whole JSON.
	for i := 1; i < len(jp.Keys); i++ {
		key := jp.Keys[i]

		switch jv.JType {
		case jutil.JTObject:
			rm := jv.Value.(map[string]*V)
			subv, ok := rm[key]
			if !ok {
				return nil, ErrNotFound
			}
			jv = subv
		case jutil.JTArray:
			ra := jv.Value.([]*V)
			if key == "-" {
				// Last element.
				if len(ra) > 0 {
					jv = ra[len(ra)-1]
				} else {
					return nil, fmt.Errorf("bad json pointer, array len is 0")
				}
			} else {
				idx, err := jutil.KeyToArrayIndex(key)
				if err != nil {
					return nil, err
				}

				if idx < 0 {
					return nil, fmt.Errorf(
						"bad json pointer, array len is %d", idx)
				}

				if idx >= len(ra) {
					return nil, ErrArrayIdxOverflow
				}
				jv = ra[idx]
			}
		default:
			return nil, fmt.Errorf(
				"bad json ponter, the value is not object or array")
		}

		if i == len(jp.Keys)-1 {
			return jv, nil
		}
	}
	return nil, fmt.Errorf("should not be here")
}

// GetArrayLen will return the size of array referenced by jp.
func (v *V) GetArrayLen(jp *ejp.ExJSONPointer) int {
	// Here, len(jp.Keys) MUST be >= 1
	if len(jp.Keys) == 1 {
		if v.JType == jutil.JTArray {
			arr := v.Value.([]*V)
			return len(arr)
		}

		return 0
	}

	jv := v
	// The first key is "", which represents the whole JSON.
	for i := 1; i < len(jp.Keys); i++ {
		key := jp.Keys[i]

		switch jv.JType {
		case jutil.JTObject:
			rm := jv.Value.(map[string]*V)
			subv, ok := rm[key]
			if !ok {
				if i == len(jp.Keys)-1 {
					return 0
				}
				return 0
			}
			jv = subv
		case jutil.JTArray:
			ra := jv.Value.([]*V)
			if key == "-" {
				// Last element.
				if len(ra) > 0 {
					jv = ra[len(ra)-1]
				} else {
					return 0
				}
			} else {
				idx, err := jutil.KeyToArrayIndex(key)
				if err != nil {
					return 0
				}

				if idx < 0 {
					return 0
				}

				if idx >= len(ra) {
					return 0
				}
				jv = ra[idx]
			}
		default:
			return 0
		}

		if i == len(jp.Keys)-1 {
			if jv.JType == jutil.JTArray {
				arr := jv.Value.([]*V)
				return len(arr)
			}
			return 0
		}
	}
	return 0
}

// Marshal ...
func (v *V) Marshal() (json.RawMessage, error) {
	obj, err := marshal(v)
	if err != nil {
		return nil, err
	}

	return json.Marshal(obj)
}

func marshal(jv *V) (interface{}, error) {
	var err error
	switch jv.JType {
	case jutil.JTObject:
		mjv := jv.Value.(map[string]*V)
		mv := make(map[string]interface{}, len(mjv))
		for k, v := range mjv {
			mv[k], err = marshal(v)
			if err != nil {
				return nil, err
			}
		}
		return mv, nil
	case jutil.JTArray:
		marr := jv.Value.([]*V)
		mv := make([]interface{}, 0, len(marr))

		var newv interface{}
		for _, arrv := range marr {
			newv, err = marshal(arrv)
			if err != nil {
				return nil, err
			}
			mv = append(mv, newv)
		}
		return mv, nil
	default:
	}

	return jv.Value, nil
}
