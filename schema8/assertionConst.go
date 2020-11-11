package schema8

import (
	"context"

	"github.com/dackon/jtool/jvalue"
)

type assertionConst struct {
	base
}

func newAssertionConst(name string, n *schemaNode) (assertion, error) {
	ast := &assertionConst{base: base{name, n}}
	return ast, nil
}

func (ast *assertionConst) Valid(ctx context.Context, v *jvalue.V) error {
	ret := isEqual(ast.node, v)
	if ret {
		return nil
	}

	return ErrAstConst
}
