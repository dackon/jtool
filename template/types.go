package template

import (
	"errors"

	"github.com/dackon/jtool/jvalue"
)

var (
	errNoValues = errors.New("no values")
)

type iterationInfo struct {
	jv         *jvalue.V
	dim        int
	maxArrSize int
	err        error
}
