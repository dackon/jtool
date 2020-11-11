package schema8

import (
	"context"
	"fmt"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type assertionFormat struct {
	base
	f      FormatFunc
	format string
}

func newAssertionFormat(name string, n *schemaNode) (assertion, error) {
	if n.jtType != jutil.JTString {
		return nil, errWithPath(
			fmt.Errorf(
				"value of format must be a string"), n)
	}

	v := n.value.(string)

	f, ok := gFormatFuncMap[v]
	if !ok || f == nil {
		return nil, errWithPath(fmt.Errorf("This implementation does not "+
			"suppoprt format '%s'", v), n)
	}

	ast := &assertionFormat{
		base:   base{name, n},
		f:      f,
		format: v,
	}
	return ast, nil
}

func (ast *assertionFormat) Valid(ctx context.Context, v *jvalue.V) error {
	err := ast.f(v)
	if err == nil {
		return nil
	}

	return fmt.Errorf("match format '%s' failed. Err is %s", ast.format, err)
}
