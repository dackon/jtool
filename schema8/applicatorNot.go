package schema8

import (
	"context"
	"fmt"

	"github.com/dackon/jtool/jvalue"
)

type applicatorNot struct {
	base
}

func newApplicatorNot(name string, n *schemaNode) (applicator, error) {
	apptor := &applicatorNot{
		base: base{name, n},
	}

	return apptor, nil
}

func (apt *applicatorNot) Valid(ctx context.Context, v *jvalue.V) error {
	info := apt.node.match(ctx, v)
	if info.Err == nil {
		err := fmt.Errorf("applicator 'not' failed")
		setFailedMatchInfo(apt.node, v, err, ctx)
		return err
	} else {
		info.Err = nil
	}

	return nil
}

func (apt *applicatorNot) Annotation() interface{} {
	return nil
}
