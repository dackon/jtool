package schema8

import (
	"context"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type applicatorPropertyNames struct {
	base
}

func newApplicatorPropertyNames(name string, n *schemaNode) (
	applicator, error) {
	apptor := &applicatorPropertyNames{
		base: base{name, n},
	}

	return apptor, nil
}

func (apt *applicatorPropertyNames) Valid(ctx context.Context,
	v *jvalue.V) error {

	if v.JType != jutil.JTObject {
		return nil
	}

	mv := v.Value.(map[string]*jvalue.V)
	for k := range mv {
		jv := &jvalue.V{JType: jutil.JTString, Value: k, Parent: v,
			KeyInParent: k}
		ret := apt.node.match(ctx, jv)
		if ret.Err != nil {
			return ret.Err
		}
	}

	return nil
}

func (apt *applicatorPropertyNames) Annotation() interface{} {
	return nil
}
