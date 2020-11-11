package schema8

import (
	"context"
	"fmt"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type assertionUniqueItems struct {
	base
	unique bool
}

func newAssertionUniqueItems(name string, n *schemaNode) (assertion, error) {
	if n.jtType != jutil.JTBoolean {
		return nil, errWithPath(
			fmt.Errorf(
				"value of uniqueItems must be a bool"), n)
	}

	v := n.value.(bool)
	ast := &assertionUniqueItems{
		base:   base{name, n},
		unique: v,
	}
	return ast, nil
}

func (ast *assertionUniqueItems) Valid(ctx context.Context, v *jvalue.V) error {
	if v.JType != jutil.JTArray || !ast.unique {
		return nil
	}

	arr := v.Value.([]*jvalue.V)
	for i := 0; i < len(arr); i++ {
		for j := i + 1; j < len(arr); j++ {
			if arr[i].IsEqual(arr[j]) == nil {
				return fmt.Errorf("Expect a unique array but the items at %d "+
					"and %d are the same", i, j)
			}
		}
	}

	return nil
}
