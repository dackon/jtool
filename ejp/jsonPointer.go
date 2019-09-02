package ejp

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// ExJSONPointer ...
type ExJSONPointer struct {
	Name string

	// Key ~0 -> ~
	// Key ~1 -> /
	// Key ~2 -> $
	Keys []string

	// Key ~ -> ~0
	// Key / -> ~1
	// Key $ -> ~2
	EncodedKeys []string
}

func (jp *ExJSONPointer) String() string {
	if len(jp.Keys) == 0 {
		return jp.Name
	}

	str := jp.Name
	for _, k := range jp.EncodedKeys {
		if str == "" {
			str = fmt.Sprintf("/%s", k)
		} else {
			str = fmt.Sprintf("%s/%s", str, k)
		}
	}

	return str
}

// Copy ...
func (jp *ExJSONPointer) Copy() *ExJSONPointer {
	cp := &ExJSONPointer{
		Name:        jp.Name,
		Keys:        make([]string, len(jp.Keys)),
		EncodedKeys: make([]string, len(jp.EncodedKeys)),
	}

	for i := 0; i < len(jp.Keys); i++ {
		cp.Keys[i] = jp.Keys[i]
		cp.EncodedKeys[i] = jp.EncodedKeys[i]
	}

	return cp
}

// GenExJSONPointer ...
func GenExJSONPointer(jsonName string, keys []string) (*ExJSONPointer, error) {
	if err := isValidJSONName(jsonName); err != nil {
		return nil, err
	}

	if len(keys) == 0 {
		return ParseExJSONPointer(jsonName)
	}

	jp := jsonName
	for i := 0; i < len(keys); i++ {
		k := keys[i]
		k = encodeJSONKey(k)
		if jp == "" {
			jp = fmt.Sprintf("/%s", k)
		} else {
			jp = fmt.Sprintf("%s/%s", jp, k)
		}
	}

	return ParseExJSONPointer(jp)
}

// ParseExJSONPointer ...
// For 'jsonFoo/a/$/~2'
//    Name: jsonFoo
//    Keys: ["", "a", "$", "$"]
//    EncodedKeys: ["", "a", "$", "~2"]
func ParseExJSONPointer(jp string) (*ExJSONPointer, error) {
	var err error

	strArr := strings.Split(jp, "/")

	// The first elem in the array must be the json name.
	if err = isValidJSONName(strArr[0]); err != nil {
		return nil, err
	}

	jpObj := &ExJSONPointer{
		Name:        strArr[0],
		EncodedKeys: make([]string, len(strArr)),
	}

	// The 1st elem in the array must be the name of the json, since we have
	// extracted the name, we can set it to empty, and use it as the root key of
	// the json.
	strArr[0] = ""

	jpObj.Keys = strArr
	for i := 0; i < len(jpObj.Keys); i++ {
		jpObj.EncodedKeys[i] = jpObj.Keys[i]
		jpObj.Keys[i], err = decodeJSONKey(jpObj.Keys[i])
		if err != nil {
			return nil, err
		}
	}

	return jpObj, nil
}

// KeyToArrayIndex ...
func KeyToArrayIndex(s string) (int, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	return int(i), err
}

// valid json name must be in [_a-zA-A0-9], or empty string.
func isValidJSONName(name string) error {
	for i := 0; i < len(name); i++ {
		v := name[i]
		if (v >= 'a' && v <= 'z') || (v >= 'A' && v <= 'Z') ||
			(v >= '0' && v <= '9') || (v == '_') || (v == '.') {
			continue
		} else {
			return fmt.Errorf("char in JSONName must be in [._a-zA-Z0-9]")
		}
	}

	return nil
}

// encodeJSONKey ...
//     ~ -> ~0
//     / -> ~1
//     $ -> ~2
func encodeJSONKey(k string) string {
	buff := bytes.NewBuffer(nil)
	for i := 0; i < len(k); i++ {
		v := k[i]
		switch v {
		case '~':
			buff.WriteByte('~')
			buff.WriteByte('0')
		case '/':
			buff.WriteByte('~')
			buff.WriteByte('1')
		case '$':
			buff.WriteByte('~')
			buff.WriteByte('2')
		default:
			buff.WriteByte(v)
		}
	}
	return buff.String()
}

// decodeJSONKey ...
//     ~0 -> ~
//     ~1 -> /
//     ~2 -> $
func decodeJSONKey(k string) (string, error) {
	buff := bytes.NewBuffer(nil)

	for i := 0; i < len(k); {
		if k[i] == '~' && i < len(k)-1 {
			switch k[i+1] {
			case '0':
				buff.WriteByte('~')
			case '1':
				buff.WriteByte('/')
			case '2':
				buff.WriteByte('$')
			default:
				return "", fmt.Errorf("bad key %s", k)
			}
			i += 2
		} else {
			buff.WriteByte(k[i])
			i++
		}
	}

	return buff.String(), nil
}
