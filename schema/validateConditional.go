package schema

import (
	"errors"

	"github.com/dackon/jtool/jvalue"
)

type condValidator struct {
	schemaImpl *schemaImpl

	// Conditional validator
	needConditionalCheck bool
	cdlIf                *Schema
	cdlThen              *Schema
	cdlElse              *Schema
}

func (cv *condValidator) loadValidator() error {
	ifv, ok := cv.schemaImpl.jvMap["if"]
	if !ok {
		cv.needConditionalCheck = false
		return nil
	}

	var err error
	thenv, ok := cv.schemaImpl.jvMap["then"]
	if !ok {
		return errors.New("'then' is absent")
	}

	elsev, ok := cv.schemaImpl.jvMap["else"]
	if !ok {
		return errors.New("'else' is absent")
	}

	key := schemaKey(cv.schemaImpl.Schema, "if")
	cv.cdlIf, err = parse(key, ifv, cv.schemaImpl.Schema,
		cv.schemaImpl.Schema.root)
	if err != nil {
		return err
	}

	key = schemaKey(cv.schemaImpl.Schema, "then")
	cv.cdlThen, err = parse(key, thenv, cv.schemaImpl.Schema,
		cv.schemaImpl.Schema.root)
	if err != nil {
		return err
	}

	key = schemaKey(cv.schemaImpl.Schema, "else")
	cv.cdlElse, err = parse(key, elsev, cv.schemaImpl.Schema,
		cv.schemaImpl.Schema.root)
	if err != nil {
		return err
	}
	cv.schemaImpl.svMap["if"] = newSchemaVSma(cv.cdlIf)
	cv.schemaImpl.svMap["then"] = newSchemaVSma(cv.cdlThen)
	cv.schemaImpl.svMap["else"] = newSchemaVSma(cv.cdlElse)
	cv.needConditionalCheck = true
	return nil
}

func (cv *condValidator) doValidate(jv *jvalue.V) (err error) {
	if !cv.needConditionalCheck {
		return nil
	}

	if err = cv.cdlIf.MatchJValue(jv); err != nil {
		return cv.cdlElse.MatchJValue(jv)
	}

	return cv.cdlThen.MatchJValue(jv)
}
