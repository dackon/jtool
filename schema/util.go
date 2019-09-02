package schema

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

const (

	// Single Precision	1 [31]	8  [30–23]	23 [22–00]
	// Double Precision	1 [63]	11 [62–52]	52 [51–00]
	gMinUint53 = 0
	gMaxUint53 = 4503599627370495  // 2^52
	gMinInt53  = -2251799813685248 // -2^51
	gMaxInt53  = 2251799813685247  // 2^51 - 1
)

func jsonPointerKeyTransform(s string) string {
	rs := bytes.NewBuffer(make([]byte, 0, len(s)))

	for i := 0; i < len(s); {
		j := i + 1
		if j >= len(s) {
			rs.WriteByte(s[i])
			return rs.String()
		}
		if s[i] == '~' && s[j] == '0' {
			rs.WriteByte('~')
			i += 2
		} else if s[i] == '~' && s[j] == '1' {
			rs.WriteByte('/')
			i += 2
		} else {
			rs.WriteByte(s[i])
			i++
		}
	}

	return rs.String()
}

func decodeJSONPointer(jp string) ([]string, error) {
	var err error

	if jp == "" {
		return []string{}, nil
	}

	if jp[0] == '#' {
		jp, err = url.PathUnescape(jp)
		if err != nil {
			return nil, err
		}

		jp = jp[1:]
	}

	if jp[0] != '/' {
		return nil, errors.New("bad json pointer")
	}

	if jp == "/" {
		return []string{""}, nil
	}

	arr := strings.Split(jp, "/")
	keys := make([]string, 0, len(arr)-1)

	for i := 1; i < len(arr); i++ {
		keys = append(keys, jsonPointerKeyTransform(arr[i]))
	}

	return keys, nil
}

func jsonKeyEncode(key string) string {
	rs := bytes.NewBuffer(make([]byte, 0, len(key)+2))
	for i := 0; i < len(key); i++ {
		switch key[i] {
		case '~':
			rs.WriteByte('~')
			rs.WriteByte('0')
		case '/':
			rs.WriteByte('~')
			rs.WriteByte('1')
		default:
			rs.WriteByte(key[i])
		}
	}

	return rs.String()
}

func schemaKey(s *Schema, key string) string {
	if s.root == s {
		return fmt.Sprintf("/%s", jsonKeyEncode(key))
	}
	return fmt.Sprintf("%s/%s", s.key, jsonKeyEncode(key))
}

func objectKey(s *Schema, key, objK string) string {
	if s.root == s {
		return fmt.Sprintf("/%s/%s", jsonKeyEncode(key), jsonKeyEncode(objK))
	}
	return fmt.Sprintf("%s/%s/%s", s.key, jsonKeyEncode(key),
		jsonKeyEncode(objK))
}

func arrayKey(s *Schema, key string, index int) string {
	if s.root == s {
		return fmt.Sprintf("/%s/%d", jsonKeyEncode(key), index)
	}
	return fmt.Sprintf("%s/%s/%d", s.key, jsonKeyEncode(key), index)
}
