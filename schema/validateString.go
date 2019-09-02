package schema

import (
	"fmt"
	"regexp"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type stringValidator struct {
	schemaImpl *schemaImpl

	// For strings.
	needMaxLengthCheck bool
	maxLength          int64
	needMinLengthCheck bool
	minLength          int64
	needPatternCheck   bool
	pattern            string
	patternRegexp      *regexp.Regexp
}

func (sv *stringValidator) loadValidator() error {
	var err error
	if err = sv.loadMaxLength(); err != nil {
		return err
	}
	if err = sv.loadMinLength(); err != nil {
		return err
	}
	if err = sv.loadPatternRegexp(); err != nil {
		return err
	}
	return nil
}

func (sv *stringValidator) doValidate(jv *jvalue.V) (err error) {
	if jv.JType != jutil.JTString {
		return nil
	}

	str, _ := jv.GetString()

	if err = sv.vMaxLength(str); err != nil {
		return err
	}

	if err = sv.vMinLength(str); err != nil {
		return err
	}

	if err = sv.vPattern(str); err != nil {
		return err
	}

	return nil
}

func (sv *stringValidator) loadMaxLength() error {
	v, ok := sv.schemaImpl.jvMap["maxLength"]
	if !ok {
		sv.needMaxLengthCheck = false
		return nil
	}

	if v.JType != jutil.JTInteger {
		sv.needMaxLengthCheck = false
		return fmt.Errorf("string maxLength must be integer")
	}

	sv.maxLength, _ = v.GetInteger()
	if sv.maxLength < 0 {
		return fmt.Errorf("maxLength must be >= 0")
	}
	sv.needMaxLengthCheck = true
	return nil
}

func (sv *stringValidator) loadMinLength() error {
	v, ok := sv.schemaImpl.jvMap["minLength"]
	if !ok {
		sv.needMinLengthCheck = false
		return nil
	}

	if v.JType != jutil.JTInteger {
		sv.needMinLengthCheck = false
		return fmt.Errorf("string minLength must be integer")
	}

	sv.minLength, _ = v.GetInteger()
	if sv.minLength < 0 {
		return fmt.Errorf("minLength must be >= 0")
	}
	sv.needMinLengthCheck = true
	return nil
}

func (sv *stringValidator) loadPatternRegexp() error {
	v, ok := sv.schemaImpl.jvMap["pattern"]
	if !ok {
		sv.needPatternCheck = false
		return nil
	}

	if v.JType != jutil.JTString {
		sv.needPatternCheck = false
		return fmt.Errorf("string pattern must be string")
	}

	sv.pattern, _ = v.GetString()

	var err error
	sv.patternRegexp, err = regexp.Compile(sv.pattern)
	if err != nil {
		return err
	}

	sv.needPatternCheck = true
	return nil
}

func (sv *stringValidator) vMaxLength(str string) (err error) {
	if !sv.needMaxLengthCheck {
		return nil
	}

	if len(str) > int(sv.maxLength) {
		return fmt.Errorf("maxLength failed. schema is %s",
			sv.schemaImpl.Schema.key)
	}

	return nil
}

func (sv *stringValidator) vMinLength(str string) (err error) {
	if !sv.needMinLengthCheck {
		return nil
	}

	if len(str) < int(sv.minLength) {
		return fmt.Errorf("minLength failed. schema is %s",
			sv.schemaImpl.Schema.key)
	}

	return nil
}

func (sv *stringValidator) vPattern(str string) (err error) {
	if !sv.needPatternCheck {
		return nil
	}

	if !sv.patternRegexp.MatchString(str) {
		return fmt.Errorf("pattern failed. schema is %s",
			sv.schemaImpl.Schema.key)
	}

	return nil
}
