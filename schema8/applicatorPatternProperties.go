package schema8

import (
	"context"
	"regexp"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type patternSchema struct {
	reg  *regexp.Regexp
	node *schemaNode
}

type applicatorPatternProperties struct {
	base
	psArr []*patternSchema
}

func newApplicatorPatternProperties(name string, n *schemaNode) (
	applicator, error) {

	apptor := &applicatorPatternProperties{
		base: base{name, n},
	}

	for k, v := range n.kvMap {
		r, err := regexp.Compile(k)
		if err != nil {
			return nil, errWithPath(err, n)
		}
		apptor.psArr = append(apptor.psArr, &patternSchema{reg: r, node: v})
	}

	return apptor, nil
}

func (apt *applicatorPatternProperties) Valid(
	ctx context.Context, v *jvalue.V) error {

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
		for _, ps := range apt.psArr {
			if ps.reg.MatchString(k) {
				ret := ps.node.match(ctx, v)
				if ret.Err != nil {
					return ret.Err
				}
			}
		}
	}

	return nil
}

func (apt *applicatorPatternProperties) Annotation() interface{} {
	return nil
}
