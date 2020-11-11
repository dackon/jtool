package schema8

import (
	"context"
	"fmt"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type assertionDependentRequired struct {
	base
	dm map[string][]string
}

func newAssertionDependentRequired(name string, n *schemaNode) (
	assertion, error) {

	if n.jtType != jutil.JTObject {
		return nil, errWithPath(fmt.Errorf(
			"value of dependentRequired must be an object"), n)
	}

	adr := &assertionDependentRequired{
		base: base{name, n},
		dm:   make(map[string][]string),
	}

	for k, v := range n.kvMap {
		if v.jtType != jutil.JTArray {
			return nil, errWithPath(fmt.Errorf(
				"each field of object must be unique string array"), v)
		}

		m := make(map[string]struct{}, len(v.nodeArr))
		arr := make([]string, 0, len(v.nodeArr))
		for _, e := range v.nodeArr {
			if e.jtType != jutil.JTString {
				return nil, errWithPath(fmt.Errorf(
					"each field of object must be unique string array"), v)

			}

			s := e.value.(string)
			_, ok := m[s]
			if ok {
				return nil, errWithPath(fmt.Errorf(
					"each field of object must be unique string array"), v)
			}

			m[s] = struct{}{}
			arr = append(arr, s)
		}
		adr.dm[k] = arr
	}

	return adr, nil
}

func (ast *assertionDependentRequired) Valid(ctx context.Context,
	v *jvalue.V) error {

	if v.JType != jutil.JTObject {
		return nil
	}

	m, _ := v.Value.(map[string]*jvalue.V)
	for k, arr := range ast.dm {
		_, ok := m[k]
		if !ok {
			continue
		}

		for _, s := range arr {
			_, ok = m[s]
			if !ok {
				return fmt.Errorf("required property '%s' "+
					"not found. Dependent key is %s", s, k)
			}
		}
	}

	return nil
}
