package schema8

import (
	"context"
	"fmt"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type assertionType struct {
	base
	TypeArr []jutil.JType
}

func newAssertionType(name string, n *schemaNode) (assertion, error) {
	at := &assertionType{
		base: base{name, n},
	}

	addType := func(s string) error {
		jt, err := jutil.IsValidJSONType(s)
		if err != nil {
			return errWithPath(err, n)
		}
		at.TypeArr = append(at.TypeArr, jt)
		return nil
	}

	switch n.jtType {
	case jutil.JTString:
		s := n.value.(string)
		if err := addType(s); err != nil {
			return nil, err
		}
		return at, nil
	case jutil.JTArray:
		for _, v := range n.nodeArr {
			if v.jtType != jutil.JTString {
				return nil, errWithPath(ErrBadType, v)
			}

			s := v.value.(string)
			if err := addType(s); err != nil {
				return nil, err
			}
		}
		return at, nil
	}

	return nil, errWithPath(ErrBadType, n)
}

func (ast *assertionType) Valid(ctx context.Context, v *jvalue.V) error {
	for _, e := range ast.TypeArr {
		if v.JType == e {
			return nil
		}

		if v.JType == jutil.JTInteger && e == jutil.JTNumber {
			return nil
		}
	}

	return fmt.Errorf("type '%s' is not in type array '%v'",
		v.JType, ast.TypeArr)
}
