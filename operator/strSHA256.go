package operator

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type strSHA256 struct {
}

func (op *strSHA256) do(jv *jvalue.V) (*jvalue.V, error) {
	if jv.JType != jutil.JTString {
		return nil, fmt.Errorf("operator 'toLower': value is not string")
	}

	s, _ := jv.GetString()
	v := sha256.Sum256([]byte(s))
	s = hex.EncodeToString(v[:])
	return &jvalue.V{JType: jutil.JTString, Value: s}, nil
}
