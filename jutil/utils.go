package jutil

import (
	"errors"
	"strconv"
)

// Errors.
var (
	ErrBadType = errors.New("unknown json type")
)

// IsSpace ...
func IsSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\r' || c == '\n'
}

// IsValidJSONType ...
func IsValidJSONType(t string) (JType, error) {
	for i := 0; i < len(GJSONTypeArr); i++ {
		if t == GJSONTypeArr[i].vstring {
			return GJSONTypeArr[i].venum, nil
		}
	}
	return JTNone, ErrBadType
}

// KeyToArrayIndex ...
func KeyToArrayIndex(s string) (int, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	return int(i), err
}

// FloatEquals ...
func FloatEquals(a, b float64) bool {
	return (a-b) < EPSILON && (b-a) < EPSILON
}
