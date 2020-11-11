package schema8

import (
	"context"

	"github.com/dackon/jtool/jvalue"
)

type applicatorIf struct {
	base
	elseApt applicator
	thenApt applicator
}

func newApplicatorIf(name string, n *schemaNode) (applicator, error) {
	apptor := &applicatorIf{
		base: base{name, n},
	}

	return apptor, nil
}

func (apt *applicatorIf) Valid(ctx context.Context, v *jvalue.V) error {
	info := apt.node.match(ctx, v)
	if info.Err == nil {
		if apt.thenApt != nil {
			return apt.thenApt.Valid(ctx, v)
		}
		return nil
	}

	// Here, 'if' schema failed, check 'else' schema.

	info.Err = nil
	if apt.elseApt != nil {
		return apt.elseApt.Valid(ctx, v)
	}

	return nil
}

func (apt *applicatorIf) Annotation() interface{} {
	return nil
}
