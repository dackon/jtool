package schema8

import (
	"context"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type applicatorItems struct {
	base
}

func newApplicatorItems(name string, n *schemaNode) (applicator, error) {
	apptor := &applicatorItems{
		base: base{name, n},
	}

	return apptor, nil
}

func (apt *applicatorItems) Valid(ctx context.Context, v *jvalue.V) error {
	if v.JType != jutil.JTArray {
		return nil
	}

	varr := v.Value.([]*jvalue.V)

	if apt.node.jtType == jutil.JTArray {
		return apt.checkItemsArray(varr, ctx)
	}

	// If here, schema 'items' is an object

	for _, e := range varr {
		ret := apt.node.match(ctx, e)
		if ret.Err != nil {
			return ret.Err
		}
	}
	return nil
}

func (apt *applicatorItems) checkItemsArray(varr []*jvalue.V,
	ctx context.Context) error {

	mc := ctx.Value(matchContextKey).(*matchContext)
	copyPath := mc.path[:]
	mc.path = append(mc.path, newMatchPathNode(apt.node))
	defer func() {
		mc.path = copyPath
	}()

	var i int
	var s *schemaNode

	for i, s = range apt.node.nodeArr {
		if i < len(varr) {
			ret := s.match(ctx, varr[i])
			if ret.Err != nil {
				return ret.Err
			}
		} else {
			return nil
		}
	}

	return nil
}

func (apt *applicatorItems) Annotation() interface{} {
	return nil
}
