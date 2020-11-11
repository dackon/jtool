package schema8

import (
	"context"
	"fmt"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type assertionEnum struct {
	base
	enum []*schemaNode
}

func newAssertionEnum(name string, n *schemaNode) (assertion, error) {
	if n.jtType != jutil.JTArray {
		return nil, errWithPath(fmt.Errorf("enum must be array"), n)
	}

	ast := &assertionEnum{
		base: base{name, n},
		enum: n.nodeArr,
	}
	return ast, nil
}

func (ast *assertionEnum) Valid(ctx context.Context, v *jvalue.V) error {
	for _, e := range ast.enum {
		if isEqual(e, v) {
			return nil
		}
	}

	return ErrAstEnum
}
