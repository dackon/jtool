package schema8

import (
	"context"

	"github.com/dackon/jtool/jvalue"
)

type applicatorThen struct {
	base
}

func newApplicatorThen(name string, n *schemaNode) (applicator, error) {
	apptor := &applicatorThen{
		base: base{name, n},
	}

	return apptor, nil
}

func (apt *applicatorThen) Valid(ctx context.Context, v *jvalue.V) error {
	return apt.node.match(ctx, v).Err
}

func (apt *applicatorThen) Annotation() interface{} {
	return nil
}
