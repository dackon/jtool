package operator

import (
	"fmt"
	"strings"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type toLower struct {
}

func (op *toLower) do(jv *jvalue.V) (*jvalue.V, error) {
	if jv.JType != jutil.JTString {
		return nil, fmt.Errorf("operator 'toLower': value is not string")
	}

	s, _ := jv.GetString()
	return &jvalue.V{JType: jutil.JTString, Value: strings.ToLower(s)}, nil
}
