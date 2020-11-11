package schema8

import (
	"context"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type applicatorProperties struct {
	base
}

func newApplicatorProperties(name string, n *schemaNode) (applicator, error) {
	apptor := &applicatorProperties{
		base{name, n},
	}

	return apptor, nil
}

func (apt *applicatorProperties) Valid(ctx context.Context, v *jvalue.V) error {
	if v.JType != jutil.JTObject {
		return nil
	}

	mc := ctx.Value(matchContextKey).(*matchContext)
	copyPath := mc.path[:]
	mc.path = append(mc.path, newMatchPathNode(apt.node))
	defer func() {
		mc.path = copyPath
	}()

	mv := v.Value.(map[string]*jvalue.V)

	for k, v := range mv {
		s, ok := apt.node.kvMap[k]
		if ok {
			ret := s.match(ctx, v)
			if ret.Err != nil {
				return ret.Err
			}
		}
	}

	return nil
}

func (apt *applicatorProperties) Annotation() interface{} {
	return nil
}
