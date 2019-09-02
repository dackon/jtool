package jutil

// JType ...
type JType int

// 7 primitive types.
const (
	JTNone JType = iota
	JTNull
	JTBoolean
	JTObject
	JTArray
	JTNumber
	JTString
	JTInteger
)

// GJSONTypeArr ...
var GJSONTypeArr = []struct {
	vstring string
	venum   JType
}{
	{"null", JTNull},
	{"boolean", JTBoolean},
	{"object", JTObject},
	{"array", JTArray},
	{"number", JTNumber},
	{"string", JTString},
	{"integer", JTInteger},
}

// Constants.
const (
	EPSILON = 0.00000001
)
