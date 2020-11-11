package schema8

import (
	"errors"
	"regexp"
)

const (
	sntBoolSchema   schemaNodeType = "BOOL_SCHEMA"
	sntObjectSchema schemaNodeType = "OBJECT_SCHEMA"
	sntNotSchema    schemaNodeType = "NOT_SCHEMA"
)

const (
	pntString  pathNodeType = "string"
	pntInteger pathNodeType = "integer"
)

const (
	matchContextKey contextKeyType = "matchContext"
)

var (
	ErrNotSchema = errors.New("input is not schema")
	ErrBadSchema = errors.New("input is bad schema")
	ErrBadType   = errors.New("type must be string or string array")
	ErrNotFound  = errors.New("schema not found")
	ErrAstConst  = errors.New("assertion 'const' match failed")
	ErrAstEnum   = errors.New("assertion 'enum' match failed")
)

var (
	gAnchorReg               = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9\-_:\.]*`)
	gAssertionConstructorMap map[string]assertionConstructorFunc

	gApplicatorArraySchemaMap  map[string]applicatorConstructorFunc
	gApplicatorObjectSchemaMap map[string]applicatorConstructorFunc
	gApplicatorSchemaMap       map[string]applicatorConstructorFunc

	gFormatFuncMap map[string]FormatFunc

	gResolverFunc ResolverFunc
)
