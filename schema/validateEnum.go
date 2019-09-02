package schema

import (
	"fmt"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type enumValidator struct {
	schemaImpl *schemaImpl

	needEnumCheck bool
	Enum          []*jvalue.V
}

func (ev *enumValidator) loadValidator() error {
	v, ok := ev.schemaImpl.jvMap["enum"]
	if !ok {
		ev.needEnumCheck = false
		return nil
	}

	if v.JType != jutil.JTArray {
		return fmt.Errorf("enum must be an array. schema is %s",
			ev.schemaImpl.Schema.key)
	}

	ev.Enum = v.Value.([]*jvalue.V)

	for i := 0; i < len(ev.Enum); i++ {
		for j := i + 1; j < len(ev.Enum); j++ {
			if err := ev.Enum[i].IsEqual(ev.Enum[j]); err == nil {
				return fmt.Errorf("enum must not be equal. schema is %s",
					ev.schemaImpl.Schema.key)
			}
		}
	}

	ev.needEnumCheck = true
	return nil
}

func (ev *enumValidator) doValidate(jv *jvalue.V) error {
	if !ev.needEnumCheck {
		return nil
	}

	for i := 0; i < len(ev.Enum); i++ {
		if err := ev.Enum[i].IsEqual(jv); err == nil {
			return nil
		}
	}

	return fmt.Errorf("match enum failed. schema is %s",
		ev.schemaImpl.Schema.key)
}
