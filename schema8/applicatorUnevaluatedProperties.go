package schema8

import (
	"context"
	"fmt"

	"github.com/dackon/jtool/jvalue"
)

type applicatorUnevaluatedProperties struct {
	base
}

func newApplicatorUnevaluatedProperties(name string, n *schemaNode) (
	applicator, error) {

	return nil, errWithPath(fmt.Errorf(
		"unevaluatedProperties is not supported"), n)
}

func (apt *applicatorUnevaluatedProperties) Valid(v *jvalue.V,
	ctx context.Context) bool {
	panic("must not be here")
}

func (apt *applicatorUnevaluatedProperties) Annotation() interface{} {
	return nil
}
