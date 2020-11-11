package schema8

import (
	"context"
	"regexp"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type applicatorAdditionalProperties struct {
	base
	properties        []string
	patternProperties []*regexp.Regexp
}

func newApplicatorAdditionalProperties(name string, n *schemaNode) (
	applicator, error) {

	apptor := &applicatorAdditionalProperties{
		base: base{name, n},
	}

	return apptor, nil
}

func (apt *applicatorAdditionalProperties) Valid(
	ctx context.Context, v *jvalue.V) error {

	if v.JType != jutil.JTObject {
		return nil
	}

	mv := v.Value.(map[string]*jvalue.V)
L1:
	for k, value := range mv {
		for _, s := range apt.properties {
			if k == s {
				continue L1
			}
		}

		for _, r := range apt.patternProperties {
			if r.MatchString(k) {
				continue L1
			}
		}

		ret := apt.node.match(ctx, value)
		if ret.Err != nil {
			return ret.Err
		}
	}

	return nil
}

func (apt *applicatorAdditionalProperties) Annotation() interface{} {
	return nil
}
