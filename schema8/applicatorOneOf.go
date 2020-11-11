package schema8

import (
	"context"
	"fmt"

	"github.com/dackon/jtool/jvalue"
)

type applicatorOneOf struct {
	base
}

func newApplicatorOneOf(name string, n *schemaNode) (applicator, error) {
	if len(n.nodeArr) == 0 {
		return nil, errWithPath(fmt.Errorf("oneOf array is empty"), n)
	}

	apptor := &applicatorOneOf{
		base{name, n},
	}

	return apptor, nil
}

func (apt *applicatorOneOf) Valid(ctx context.Context, v *jvalue.V) error {
	mc := ctx.Value(matchContextKey).(*matchContext)
	copyPath := mc.path[:]
	mc.path = append(mc.path, newMatchPathNode(apt.node))
	defer func() {
		mc.path = copyPath
	}()

	count := 0
	for _, n := range apt.node.nodeArr {
		info := n.match(ctx, v)
		if info.Err == nil {
			count++
		} else {
			info.Err = nil
		}
	}

	if count == 1 {
		return nil
	}

	err := fmt.Errorf("applicator 'oneOf' failed. %d elements matched", count)
	setFailedMatchInfo(apt.node, v, err, ctx)
	return err
}

func (apt *applicatorOneOf) Annotation() interface{} {
	return nil
}
