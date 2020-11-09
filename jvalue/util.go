package jvalue

import (
	"bytes"
)

func DecodeJSONKey(s string) string {
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

func EncodeJSONKey(key string) string {
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
