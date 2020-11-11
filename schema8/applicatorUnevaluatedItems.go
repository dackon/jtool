package schema8

import (
	"context"
	"fmt"

	"github.com/dackon/jtool/jvalue"
)

type applicatorUnevaluatedItems struct {
	base
}

func newApplicatorUnevaluatedItems(name string, n *schemaNode) (
	applicator, error) {
	return nil, errWithPath(fmt.Errorf("unevaluatedItems is not supported"), n)
}

func (apt *applicatorUnevaluatedItems) Valid(v *jvalue.V,
	ctx context.Context) bool {
	panic("must not be here")
}

func (apt *applicatorUnevaluatedItems) Annotation() interface{} {
	return nil
}
