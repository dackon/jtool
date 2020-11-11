package schema8

import (
	"context"
	"fmt"

	"github.com/dackon/jtool/jvalue"
)

type applicatorAnyOf struct {
	base
}

func newApplicatorAnyOf(name string, n *schemaNode) (applicator, error) {
	if len(n.nodeArr) == 0 {
		return nil, errWithPath(fmt.Errorf("allOf array is empty"), n)
	}

	apptor := &applicatorAnyOf{base{name, n}}
	return apptor, nil
}

func (apt *applicatorAnyOf) Valid(ctx context.Context, v *jvalue.V) error {
	mc := ctx.Value(matchContextKey).(*matchContext)
	copyPath := mc.path[:]
	mc.path = append(mc.path, newMatchPathNode(apt.node))
	defer func() {
		mc.path = copyPath
	}()

	for _, n := range apt.node.nodeArr {
		info := n.match(ctx, v)
		if info.Err == nil {
			return nil
		} else {
			info.Err = nil
		}
	}

	err := fmt.Errorf("applicator 'anyOf' failed")
	setFailedMatchInfo(apt.node, v, err, ctx)
	return err
}

func (apt *applicatorAnyOf) Annotation() interface{} {
	return nil
}
