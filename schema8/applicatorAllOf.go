package schema8

import (
	"context"
	"fmt"

	"github.com/dackon/jtool/jvalue"
)

type applicatorAllOf struct {
	base
}

func newApplicatorAllOf(name string, n *schemaNode) (applicator, error) {
	if len(n.nodeArr) == 0 {
		return nil, errWithPath(fmt.Errorf("allOf array is empty"), n)
	}

	apptor := &applicatorAllOf{base{name, n}}
	return apptor, nil
}

func (apt *applicatorAllOf) Valid(ctx context.Context, v *jvalue.V) error {
	mc := ctx.Value(matchContextKey).(*matchContext)
	copyPath := mc.path[:]
	mc.path = append(mc.path, newMatchPathNode(apt.node))
	defer func() {
		mc.path = copyPath
	}()

	for _, n := range apt.node.nodeArr {
		info := n.match(ctx, v)
		if info.Err != nil {
			return info.Err
		}
	}

	return nil
}

func (apt *applicatorAllOf) Annotation() interface{} {
	return nil
}
