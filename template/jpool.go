package template

import (
	"fmt"

	"github.com/dackon/jtool/ejp"
	"github.com/dackon/jtool/jvalue"
)

// JSONPool ...
type JSONPool struct {
	ValueMap map[string]*jvalue.V
}

// NewJSONPool ...
func NewJSONPool() *JSONPool {
	return &JSONPool{
		ValueMap: make(map[string]*jvalue.V, 8),
	}
}

// Add ...
func (jpool *JSONPool) Add(name string, v *jvalue.V) {
	jpool.ValueMap[name] = v
}

// GetJV ...
func (jpool *JSONPool) GetJV(name string) (*jvalue.V, bool) {
	v, ok := jpool.ValueMap[name]
	return v, ok
}

// GetValue ...
func (jpool *JSONPool) GetValue(jp *ejp.ExJSONPointer) (
	*jvalue.V, error) {
	jv, ok := jpool.ValueMap[jp.Name]
	if !ok {
		return nil, fmt.Errorf("failed to find JSON name %s", jp.Name)
	}
	return jv.GetSubJSONValue(jp)
}

// GetArrayLen ...
func (jpool *JSONPool) GetArrayLen(jp *ejp.ExJSONPointer) int {
	jv, ok := jpool.ValueMap[jp.Name]
	if !ok {
		return 0
	}
	return jv.GetArrayLen(jp)
}
