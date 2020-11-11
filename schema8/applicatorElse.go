package schema8

import (
	"context"

	"github.com/dackon/jtool/jvalue"
)

type applicatorElse struct {
	base
}

func newApplicatorElse(name string, n *schemaNode) (applicator, error) {
	apptor := &applicatorElse{
		base: base{name, n},
	}

	return apptor, nil
}

func (apt *applicatorElse) Valid(ctx context.Context, v *jvalue.V) error {
	return apt.node.match(ctx, v).Err
}

func (apt *applicatorElse) Annotation() interface{} {
	return nil
}
