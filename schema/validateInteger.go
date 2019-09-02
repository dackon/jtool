package schema

import (
	"fmt"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type integerValidator struct {
	schemaImpl *schemaImpl

	// For integers.
	needMultipleOfCheck       bool
	multipleOf                int64
	needMaximumCheck          bool
	maximum                   int64
	needExclusiveMaximumCheck bool
	exclusiveMaximum          int64
	needMinimumCheck          bool
	minimum                   int64
	needExclusiveMinimumCheck bool
	exclusiveMinimum          int64
}

func (nv *integerValidator) loadValidator() error {
	var err error
	if err = nv.loadMultipleOf(); err != nil {
		return err
	}

	if err = nv.loadMaximum(); err != nil {
		return err
	}

	if err = nv.loadMinimum(); err != nil {
		return err
	}

	if err = nv.loadExclusiveMaximum(); err != nil {
		return err
	}

	if err = nv.loadExclusiveMinimum(); err != nil {
		return err
	}

	return nil
}

func (nv *integerValidator) doValidate(jv *jvalue.V) (err error) {
	if jv.JType != jutil.JTInteger {
		return nil
	}

	num, _ := jv.GetInteger()

	if err = nv.vExclusiveMaximum(num); err != nil {
		return err
	}

	if err = nv.vExclusiveMinimum(num); err != nil {
		return err
	}

	if err = nv.vMaximum(num); err != nil {
		return err
	}

	if err = nv.vMinimum(num); err != nil {
		return err
	}

	if err = nv.vMultipleOf(num); err != nil {
		return err
	}
	return nil
}

func (nv *integerValidator) loadMultipleOf() error {
	v, ok := nv.schemaImpl.jvMap["multipleOf"]
	if !ok {
		nv.needMultipleOfCheck = false
		return nil
	}

	if v.JType != jutil.JTInteger {
		// Integer is a special type, multipleOf can be used for type 'number'
		// too, so we need to check the schema type. If schema is not integer,
		// we should throw an error.
		if nv.schemaImpl.isType(jutil.JTInteger) {
			return fmt.Errorf("integer multipleOf must be integer")
		}
		nv.needMaximumCheck = false
		return nil
	}

	nv.multipleOf, _ = v.GetInteger()
	if nv.multipleOf <= 0 {
		return fmt.Errorf("multipleOf must be > 0")
	}

	nv.needMultipleOfCheck = true
	return nil
}

func (nv *integerValidator) loadMaximum() error {
	v, ok := nv.schemaImpl.jvMap["maximum"]
	if !ok {
		nv.needMaximumCheck = false
		return nil
	}

	if v.JType != jutil.JTInteger {
		if nv.schemaImpl.isType(jutil.JTInteger) {
			return fmt.Errorf("integer maximum must be integer")
		}

		nv.needMaximumCheck = false
		return nil
	}

	nv.maximum, _ = v.GetInteger()
	nv.needMaximumCheck = true
	return nil
}

func (nv *integerValidator) loadExclusiveMaximum() error {
	v, ok := nv.schemaImpl.jvMap["exclusiveMaximum"]
	if !ok {
		nv.needExclusiveMaximumCheck = false
		return nil
	}

	if v.JType != jutil.JTInteger {
		if nv.schemaImpl.isType(jutil.JTInteger) {
			return fmt.Errorf("integer exclusiveMaximum must be integer")
		}

		nv.needExclusiveMaximumCheck = false
		return nil
	}

	nv.exclusiveMaximum, _ = v.GetInteger()
	nv.needExclusiveMaximumCheck = true
	return nil
}

func (nv *integerValidator) loadMinimum() error {
	v, ok := nv.schemaImpl.jvMap["minimum"]
	if !ok {
		nv.needMinimumCheck = false
		return nil
	}

	if v.JType != jutil.JTInteger {
		if nv.schemaImpl.isType(jutil.JTInteger) {
			return fmt.Errorf("integer minimum must be integer")
		}

		nv.needMinimumCheck = false
		return nil
	}

	nv.minimum, _ = v.GetInteger()
	nv.needMinimumCheck = true
	return nil
}

func (nv *integerValidator) loadExclusiveMinimum() error {
	v, ok := nv.schemaImpl.jvMap["exclusiveMinimum"]
	if !ok {
		nv.needExclusiveMinimumCheck = false
		return nil
	}

	if v.JType != jutil.JTInteger {
		if nv.schemaImpl.isType(jutil.JTInteger) {
			return fmt.Errorf("integer exclusiveMinimum must be integer")
		}

		nv.needExclusiveMinimumCheck = false
		return nil
	}

	nv.exclusiveMinimum, _ = v.GetInteger()
	nv.needExclusiveMinimumCheck = true
	return nil
}

func (nv *integerValidator) vMultipleOf(num int64) (err error) {
	if !nv.needMultipleOfCheck {
		return nil
	}

	if num >= nv.multipleOf && num%nv.multipleOf == 0 {
		return nil
	}

	return fmt.Errorf("multipleOf failed. schema is %s",
		nv.schemaImpl.Schema.key)
}

func (nv *integerValidator) vMaximum(num int64) (err error) {
	if !nv.needMaximumCheck {
		return nil
	}

	if nv.maximum < num {
		return fmt.Errorf("maximum failed. schema is %s",
			nv.schemaImpl.Schema.key)
	}

	return nil
}

func (nv *integerValidator) vExclusiveMaximum(num int64) (err error) {
	if !nv.needExclusiveMaximumCheck {
		return nil
	}

	if nv.exclusiveMaximum <= num {
		return fmt.Errorf("exclusiveMaximum failed. schema is %s",
			nv.schemaImpl.Schema.key)
	}

	return nil
}

func (nv *integerValidator) vMinimum(num int64) (err error) {
	if !nv.needMinimumCheck {
		return nil
	}

	if nv.minimum > num {
		return fmt.Errorf("minimum failed. schema is %s",
			nv.schemaImpl.Schema.key)
	}

	return nil
}

func (nv *integerValidator) vExclusiveMinimum(num int64) (err error) {
	if !nv.needExclusiveMinimumCheck {
		return nil
	}

	if nv.exclusiveMinimum >= num {
		return fmt.Errorf("exclusiveMinimum failed. schema is %s",
			nv.schemaImpl.Schema.key)
	}

	return nil
}
