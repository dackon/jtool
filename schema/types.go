package schema

import "github.com/dackon/jtool/jvalue"

type schemaKind int

const (
	schemaBool   schemaKind = 1
	schemaObject schemaKind = 2
)

type validatorInf interface {
	loadValidator() error
	doValidate(jv *jvalue.V) error
}

// FormatFunc ...
type FormatFunc func(*jvalue.V) error
