package schema8

import (
	"context"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type applicatorAdditionalItems struct {
	base
	itemsArrLen int
}

func newApplicatorAdditionalItems(name string, n *schemaNode) (
	applicator, error) {
	apptor := &applicatorAdditionalItems{
		base: base{name, n},
	}

	return apptor, nil
}

func (apt *applicatorAdditionalItems) Valid(ctx context.Context,
	v *jvalue.V) error {

	if v.JType != jutil.JTArray {
		return nil
	}

	varr := v.Value.([]*jvalue.V)
	if len(varr) <= apt.itemsArrLen {
		return nil
	}

	for i := apt.itemsArrLen; i < len(varr); i++ {
		ret := apt.node.match(ctx, varr[i])
		if ret.Err != nil {
			return ret.Err
		}
	}

	return nil
}

func (apt *applicatorAdditionalItems) Annotation() interface{} {
	return nil
}
