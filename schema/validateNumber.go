package schema

import (
	"errors"
	"fmt"
	"math"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type numberValidator struct {
	schemaImpl *schemaImpl

	// For numbers.
	needMultipleOfCheck       bool
	multipleOf                float64
	needMaximumCheck          bool
	maximum                   float64
	needExclusiveMaximumCheck bool
	exclusiveMaximum          float64
	needMinimumCheck          bool
	minimum                   float64
	needExclusiveMinimumCheck bool
	exclusiveMinimum          float64
}

func (nv *numberValidator) loadValidator() error {
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

func (nv *numberValidator) doValidate(jv *jvalue.V) (err error) {
	if jv.JType != jutil.JTNumber && jv.JType != jutil.JTInteger {
		return nil
	}

	num, _ := jv.GetNumber()

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

func (nv *numberValidator) loadMultipleOf() error {
	v, ok := nv.schemaImpl.jvMap["multipleOf"]
	if !ok {
		nv.needMultipleOfCheck = false
		return nil
	}

	if v.JType != jutil.JTNumber && v.JType != jutil.JTInteger {
		nv.needMultipleOfCheck = false
		return fmt.Errorf("number multipleOf must be float or integer")
	}

	nv.multipleOf, _ = v.GetNumber()
	if nv.multipleOf <= 0 {
		return errors.New("multipleOf value must be " +
			"strictly greater than 0")
	}

	nv.needMultipleOfCheck = true
	return nil
}

func (nv *numberValidator) loadMaximum() error {
	v, ok := nv.schemaImpl.jvMap["maximum"]
	if !ok {
		nv.needMaximumCheck = false
		return nil
	}

	if v.JType != jutil.JTNumber && v.JType != jutil.JTInteger {
		nv.needMaximumCheck = false
		return fmt.Errorf("number maximum must be float or integer")
	}

	nv.maximum, _ = v.GetNumber()
	nv.needMaximumCheck = true
	return nil
}

func (nv *numberValidator) loadExclusiveMaximum() error {
	v, ok := nv.schemaImpl.jvMap["exclusiveMaximum"]
	if !ok {
		nv.needExclusiveMaximumCheck = false
		return nil
	}

	if v.JType != jutil.JTNumber && v.JType != jutil.JTInteger {
		nv.needExclusiveMaximumCheck = false
		return fmt.Errorf("number exclusiveMaximum must be float or integer")
	}

	nv.exclusiveMaximum, _ = v.GetNumber()
	nv.needExclusiveMaximumCheck = true
	return nil
}

func (nv *numberValidator) loadMinimum() error {
	v, ok := nv.schemaImpl.jvMap["minimum"]
	if !ok {
		nv.needMinimumCheck = false
		return nil
	}

	if v.JType != jutil.JTNumber && v.JType != jutil.JTInteger {
		nv.needMinimumCheck = false
		return fmt.Errorf("number minimum must be float or integer")
	}

	nv.minimum, _ = v.GetNumber()
	nv.needMinimumCheck = true
	return nil
}

func (nv *numberValidator) loadExclusiveMinimum() error {
	v, ok := nv.schemaImpl.jvMap["exclusiveMinimum"]
	if !ok {
		nv.needExclusiveMinimumCheck = false
		return nil
	}

	if v.JType != jutil.JTNumber && v.JType != jutil.JTInteger {
		nv.needExclusiveMinimumCheck = false
		return fmt.Errorf("number exclusiveMinimum must be float or integer")
	}

	nv.exclusiveMinimum, _ = v.GetNumber()
	nv.needExclusiveMinimumCheck = true
	return nil
}

func (nv *numberValidator) vMultipleOf(num float64) (err error) {
	if !nv.needMultipleOfCheck {
		return nil
	}

	div := num / nv.multipleOf
	ndiv := int64(div)
	fdiv := float64(ndiv)

	if math.Abs(fdiv-div) >= 1e-9 {
		return fmt.Errorf("multipleOf failed. schema is %s",
			nv.schemaImpl.Schema.key)
	}

	return nil
}

func (nv *numberValidator) vMaximum(num float64) (err error) {
	if !nv.needMaximumCheck {
		return nil
	}

	if nv.maximum < num {
		return fmt.Errorf("Maximum failed. schema is %s",
			nv.schemaImpl.Schema.key)
	}

	return nil
}

func (nv *numberValidator) vExclusiveMaximum(num float64) (err error) {
	if !nv.needExclusiveMaximumCheck {
		return nil
	}

	if nv.exclusiveMaximum <= num {
		return fmt.Errorf("exclusiveMaximum failed. schema is %s",
			nv.schemaImpl.Schema.key)
	}

	return nil
}

func (nv *numberValidator) vMinimum(num float64) (err error) {
	if !nv.needMinimumCheck {
		return nil
	}

	if nv.minimum > num {
		return fmt.Errorf("minimum failed. schema is %s",
			nv.schemaImpl.Schema.key)
	}

	return nil
}

func (nv *numberValidator) vExclusiveMinimum(num float64) (err error) {
	if !nv.needExclusiveMinimumCheck {
		return nil
	}

	if nv.exclusiveMinimum >= num {
		return fmt.Errorf("exclusiveMinimum failed. schema is %s",
			nv.schemaImpl.Schema.key)
	}

	return nil
}
