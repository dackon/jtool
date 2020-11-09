package schema

import (
	"errors"
	"fmt"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type boolLogicValidator struct {
	schemaImpl *schemaImpl

	// Boolean logic
	needAllOfCheck bool
	allOf          []*Schema
	needAnyOfCheck bool
	anyOf          []*Schema
	needOneOfCheck bool
	oneOf          []*Schema
	needNotCheck   bool
	not            *Schema
}

func (bv *boolLogicValidator) loadValidator() error {
	var err error

	if err = bv.loadAllOf(); err != nil {
		return err
	}

	if err = bv.loadAnyOf(); err != nil {
		return err
	}

	if err = bv.loadNot(); err != nil {
		return err
	}

	if err = bv.loadOneOf(); err != nil {
		return err
	}

	return nil
}

func (bv *boolLogicValidator) doValidate(jv *jvalue.V) (err error) {
	if err = bv.vAllOf(jv); err != nil {
		return err
	}

	if err = bv.vAnyOf(jv); err != nil {
		return err
	}

	if err = bv.vOneOf(jv); err != nil {
		return err
	}

	if err = bv.vNot(jv); err != nil {
		return err
	}

	return nil
}

func (bv *boolLogicValidator) loadAllOf() error {
	v, ok := bv.schemaImpl.jvMap["allOf"]
	if !ok {
		bv.needAllOfCheck = false
		return nil
	}

	if v.JType != jutil.JTArray {
		return errors.New("'allOf' is not an array")
	}

	arr := v.Value.([]*jvalue.V)
	if len(arr) == 0 {
		return errors.New("'allOf' is an empty array")
	}

	for i := 0; i < len(arr); i++ {
		key := arrayKey(bv.schemaImpl.Schema, "allOf", i)
		s, err := parse(key, arr[i], bv.schemaImpl.Schema,
			bv.schemaImpl.Schema.root)
		if err != nil {
			return err
		}

		bv.allOf = append(bv.allOf, s)
	}

	bv.schemaImpl.svMap["allOf"] = newSchemaVArr(bv.allOf)
	bv.needAllOfCheck = true
	return nil
}

func (bv *boolLogicValidator) loadAnyOf() error {
	v, ok := bv.schemaImpl.jvMap["anyOf"]
	if !ok {
		bv.needAnyOfCheck = false
		return nil
	}

	if v.JType != jutil.JTArray {
		return errors.New("'anyOf' is not an array")
	}

	arr := v.Value.([]*jvalue.V)
	if len(arr) == 0 {
		return errors.New("'anyOf' is an empty array")
	}

	for i := 0; i < len(arr); i++ {
		key := arrayKey(bv.schemaImpl.Schema, "anyOf", i)
		s, err := parse(key, arr[i], bv.schemaImpl.Schema,
			bv.schemaImpl.Schema.root)
		if err != nil {
			return err
		}

		bv.anyOf = append(bv.anyOf, s)
	}

	bv.schemaImpl.svMap["anyOf"] = newSchemaVArr(bv.anyOf)
	bv.needAnyOfCheck = true
	return nil
}

func (bv *boolLogicValidator) loadOneOf() error {
	v, ok := bv.schemaImpl.jvMap["oneOf"]
	if !ok {
		bv.needOneOfCheck = false
		return nil
	}

	if v.JType != jutil.JTArray {
		return errors.New("'oneOf' is not an array")
	}

	arr := v.Value.([]*jvalue.V)
	if len(arr) == 0 {
		return errors.New("'oneOf' is an empty array")
	}

	for i := 0; i < len(arr); i++ {
		key := arrayKey(bv.schemaImpl.Schema, "oneOf", i)
		s, err := parse(key, arr[i], bv.schemaImpl.Schema,
			bv.schemaImpl.Schema.root)
		if err != nil {
			return err
		}

		bv.oneOf = append(bv.oneOf, s)
	}

	bv.schemaImpl.svMap["oneOf"] = newSchemaVArr(bv.oneOf)
	bv.needOneOfCheck = true
	return nil
}

func (bv *boolLogicValidator) loadNot() error {
	v, ok := bv.schemaImpl.jvMap["not"]
	if !ok {
		bv.needNotCheck = false
		return nil
	}

	key := schemaKey(bv.schemaImpl.Schema, "not")
	s, err := parse(key, v, bv.schemaImpl.Schema, bv.schemaImpl.Schema.root)
	if err != nil {
		return err
	}

	bv.not = s
	bv.needNotCheck = true
	bv.schemaImpl.svMap["not"] = newSchemaVSma(bv.not)
	return nil
}

func (bv *boolLogicValidator) vAllOf(jv *jvalue.V) (err error) {
	if !bv.needAllOfCheck {
		return nil
	}

	for i := 0; i < len(bv.allOf); i++ {
		err = bv.allOf[i].MatchJValue(jv)
		if err != nil {
			return err
		}
	}

	return nil
}

func (bv *boolLogicValidator) vAnyOf(jv *jvalue.V) (err error) {
	if !bv.needAnyOfCheck {
		return nil
	}

	for i := 0; i < len(bv.anyOf); i++ {
		if err = bv.anyOf[i].MatchJValue(jv); err == nil {
			return nil
		}
	}

	return fmt.Errorf("failed to match anyOf. Schema is %s",
		bv.schemaImpl.Schema.key)
}

func (bv *boolLogicValidator) vOneOf(jv *jvalue.V) (err error) {
	if !bv.needOneOfCheck {
		return nil
	}

	count := 0
	for i := 0; i < len(bv.oneOf); i++ {
		if err = bv.oneOf[i].MatchJValue(jv); err == nil {
			count++
		}

		if count > 1 {
			return fmt.Errorf("failed to match oneOf. Schema is %s",
				bv.schemaImpl.Schema.key)
		}
	}

	if count == 1 {
		return nil
	}

	return fmt.Errorf("failed to match oneOf. Schema is %s",
		bv.schemaImpl.Schema.key)
}

func (bv *boolLogicValidator) vNot(jv *jvalue.V) (err error) {
	if !bv.needNotCheck {
		return nil
	}

	if err = bv.not.MatchJValue(jv); err == nil {
		return fmt.Errorf("failed to match not. Schema is %s",
			bv.schemaImpl.Schema.key)
	}

	return nil
}
