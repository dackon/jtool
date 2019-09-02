package schema

import "github.com/dackon/jtool/jvalue"

type formatValidator struct {
	schemaImpl *schemaImpl
}

func (fv *formatValidator) loadValidator() error {
	return nil
}

func (fv *formatValidator) doValidate(jv *jvalue.V) error {
	return nil
}
