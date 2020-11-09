package schema

import (
	"fmt"

	"github.com/dackon/jtool/jvalue"
)

type formatValidator struct {
	schemaImpl *schemaImpl

	needFormatCheck bool
	format          string
	fmtFunc         FormatFunc
}

func (fv *formatValidator) loadValidator() error {
	var err error

	v, ok := fv.schemaImpl.jvMap["format"]
	if !ok {
		fv.needFormatCheck = false
		return nil
	}

	fv.format, err = v.GetString()
	if err != nil || fv.format == "" {
		return fmt.Errorf("failed to get format. err is '%s' or format "+
			"is empty", err)
	}

	fv.fmtFunc, ok = gFMTFuncMap[fv.format]
	if !ok {
		return fmt.Errorf("format '%s' validate function not found", fv.format)
	}

	fv.needFormatCheck = true
	return nil
}

func (fv *formatValidator) doValidate(jv *jvalue.V) error {
	if !fv.needFormatCheck {
		return nil
	}

	return fv.fmtFunc(jv)
}
