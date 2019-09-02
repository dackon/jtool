package schema

import (
	"fmt"
	"strconv"
)

type svType int

const (
	svSchema svType = iota
	svObject
	svArray
)

type objectValueInterface interface {
	getSchema(key string) (*Schema, error)
}

type arrayValueInterface interface {
	getSchema(i int) (*Schema, error)
}

type schemaValue struct {
	svType svType
	value  interface{}
}

func newSchemaVObj(objv map[string]*Schema) *schemaValue {
	return &schemaValue{
		svType: svObject,
		value:  objv,
	}
}

func newSchemaVArr(arr []*Schema) *schemaValue {
	return &schemaValue{
		svType: svArray,
		value:  arr,
	}
}

func newSchemaVSma(s *Schema) *schemaValue {
	return &schemaValue{
		svType: svSchema,
		value:  s,
	}
}

func (sv *schemaValue) getSchema(key string) (*Schema, error) {
	switch sv.svType {
	case svArray:
		n, err := strconv.ParseInt(key, 10, 64)
		if err != nil {
			return nil, err
		}
		v := sv.value.([]*Schema)

		if int(n) >= len(v) {
			return nil, fmt.Errorf("failed to find schema at %d", n)
		}

		return v[n], nil
	case svObject:
		v := sv.value.(map[string]*Schema)
		s, ok := v[key]
		if !ok {
			return nil, fmt.Errorf("failed to find schema by key %s", key)
		}
		return s, nil
	case svSchema:
		return sv.value.(*Schema), nil
	default:
		panic("must not be here")
	}
}

func (sv *schemaValue) resolveRef(revArr []*Schema) error {
	return nil
}
