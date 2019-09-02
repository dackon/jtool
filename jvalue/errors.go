package jvalue

import "errors"

// Errors
var (
	ErrBadParam         = errors.New("bad parameter")
	ErrBadType          = errors.New("unknown json type")
	ErrTypeNotEqual     = errors.New("type not equal")
	ErrBoolNotEqual     = errors.New("bool not equal")
	ErrNumberNotEqual   = errors.New("number not equal")
	ErrIntegerNotEqual  = errors.New("integer not equal")
	ErrStringNotEqual   = errors.New("string not equal")
	ErrObjectNotEqual   = errors.New("object not equal")
	ErrArrayNotEqual    = errors.New("array not equal")
	ErrNotFound         = errors.New("field not exist")
	ErrArrayIdxOverflow = errors.New("bad array idx")
)
