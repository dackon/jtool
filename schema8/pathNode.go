package schema8

import (
	"fmt"

	"github.com/dackon/jtool/jvalue"
)

type pathNode struct {
	Type pathNodeType
	// If string, it is JSON pointer encoded, if integer, equals to Raw
	Value interface{}
	// Before JSON pointer encoding if it is string
	Raw interface{}
}

func newPathNodeString(s string) *pathNode {
	return &pathNode{
		Type:  pntString,
		Value: jvalue.EncodeJSONKey(s),
		Raw:   s,
	}
}

func newPathNodeInt(i int) *pathNode {
	return &pathNode{
		Type:  pntInteger,
		Value: i,
		Raw:   i,
	}
}

func (n *pathNode) String() string {
	if n.Type == pntString {
		return n.Value.(string)
	}

	return fmt.Sprintf("%d", n.Value.(int))
}
