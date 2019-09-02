package template

import (
	"github.com/dackon/jtool/ejp"
	"github.com/dackon/jtool/jutil"
)

const (
	jtPointer jutil.JType = 100
)

type tplPointer struct {
	jp             *ejp.ExJSONPointer
	wildcardIdxArr []int
}

func newTPLPointer(jp *ejp.ExJSONPointer) *tplPointer {
	tp := &tplPointer{
		jp: jp,
	}
	for i, v := range tp.jp.Keys {
		if v == "$" && tp.jp.EncodedKeys[i] == "$" {
			tp.wildcardIdxArr = append(tp.wildcardIdxArr, i)
		}
	}
	return tp
}

func (tp *tplPointer) getExJSONPointer() *ejp.ExJSONPointer {
	if len(tp.wildcardIdxArr) == 0 {
		return tp.jp
	}
	return tp.jp.Copy()
}
