package operator

import (
	"github.com/dackon/jtool/jvalue"
)

// Operator ...
type Operator interface {
	do(jv *jvalue.V) (*jvalue.V, error)
}

var (
	// Number operations.
	gAbs       = &abs{}
	gCeil      = &ceil{}
	gDivide    = &divide{}
	gSum       = &sum{}
	gDateAddMS = newDateAddMS()

	// String operations.
	gConcat  = &concat{}
	gSHA256  = &strSHA256{}
	gToLower = &toLower{}
	gToUpper = &toUpper{}

	// Array opearions.
	gLoop = &loop{}
)

var (
	gOPMap = map[string]Operator{
		"$abs":    gAbs,
		"$ceil":   gCeil,
		"$divide": gDivide,
		"$sum":    gSum,

		"$concat":  gConcat,
		"$sha256":  gSHA256,
		"$toLower": gToLower,
		"$toUpper": gToUpper,

		"$loop": gLoop,

		"$dateAddMS": gDateAddMS,
	}
)
