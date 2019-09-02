package schema

import (
	"fmt"

	"github.com/dackon/jtool/jvalue"
)

type constValidator struct {
	schemaImpl *schemaImpl

	needConstCheck bool
	constV         *jvalue.V
}

func (cv *constValidator) loadValidator() error {
	v, ok := cv.schemaImpl.jvMap["const"]
	if !ok {
		cv.needConstCheck = false
		return nil
	}

	cv.constV = v
	cv.needConstCheck = true
	return nil
}

func (cv *constValidator) doValidate(jv *jvalue.V) (err error) {
	if !cv.needConstCheck {
		return nil
	}

	err = cv.constV.IsEqual(jv)
	if err != nil {
		return fmt.Errorf("const failed. err is %s. schema is %s",
			err, cv.schemaImpl.Schema.key)
	}

	return nil
}
