package schema8

import (
	"context"
	"fmt"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type assertionRequired struct {
	base
	required []string
}

func newAssertionRequired(name string, n *schemaNode) (assertion, error) {
	if n.jtType != jutil.JTArray {
		return nil, errWithPath(fmt.Errorf(
			"value of required must be a string array"), n)
	}

	ar := &assertionRequired{
		base:     base{name, n},
		required: make([]string, 0, len(n.nodeArr)),
	}
	m := make(map[string]struct{}, len(n.nodeArr))
	for _, e := range n.nodeArr {
		if e.jtType != jutil.JTString {
			return nil, errWithPath(fmt.Errorf(
				"value of required must be a string array"), n)
		}

		s := e.value.(string)
		_, ok := m[s]
		if ok {
			return nil, errWithPath(fmt.Errorf(
				"value of required must be a unique string array"), n)
		}
		m[s] = struct{}{}
		ar.required = append(ar.required, s)
	}

	return ar, nil
}

func (ast *assertionRequired) Valid(ctx context.Context, v *jvalue.V) error {
	if v.JType != jutil.JTObject {
		return nil
	}

	m, _ := v.Value.(map[string]*jvalue.V)
	for _, r := range ast.required {
		_, ok := m[r]
		if !ok {
			return fmt.Errorf("required property '%s' not found", r)
		}
	}

	return nil
}
