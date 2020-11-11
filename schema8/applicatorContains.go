package schema8

import (
	"context"
	"fmt"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type applicatorContains struct {
	base
	minContains int
	maxContains int
}

func newApplicatorContains(name string, n *schemaNode) (applicator, error) {
	apptor := &applicatorContains{
		base:        base{name, n},
		minContains: -1,
		maxContains: -1,
	}
	return apptor, nil
}

func (apt *applicatorContains) Valid(ctx context.Context, v *jvalue.V) error {
	if v.JType != jutil.JTArray {
		return nil
	}

	var err error
	varr := v.Value.([]*jvalue.V)
	if len(varr) == 0 && apt.minContains < 0 {
		err = fmt.Errorf("applicator 'contains' failed: empty array")
		setFailedMatchInfo(apt.node, v, err, ctx)
		return err
	}

	count := 0
	for _, e := range varr {
		info := apt.node.match(ctx, e)
		if info.Err == nil {
			if apt.minContains < 0 && apt.maxContains < 0 {
				return nil
			}
			count++
		} else {
			info.Err = nil
		}
	}

	if apt.minContains < 0 && apt.maxContains < 0 {
		err = fmt.Errorf("applicator 'contains' failed")
		setFailedMatchInfo(apt.node, v, err, ctx)
		return err
	}

	if apt.minContains >= 0 && count < apt.minContains {
		err = fmt.Errorf("Expected min contains %d but found %d",
			apt.minContains, count)
		setFailedMatchInfo(apt.node, v, err, ctx)
		return err
	}

	if apt.maxContains >= 0 && count > apt.maxContains {
		err = fmt.Errorf("Expected max contains %d but found %d",
			apt.maxContains, count)
		setFailedMatchInfo(apt.node, v, err, ctx)
		return err
	}

	return nil
}

func (apt *applicatorContains) Annotation() interface{} {
	return nil
}
