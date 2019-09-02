package schema

import (
	"errors"
	"fmt"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type arrayValidator struct {
	schemaImpl *schemaImpl

	// For array.
	needItemsCheck           bool
	items                    *Schema
	needItemsArrCheck        bool
	itemsArr                 []*Schema
	needAdditionalItemsCheck bool
	additionalItems          *Schema
	needContainsCheck        bool
	contains                 *Schema
	needMaxItemsCheck        bool
	maxItems                 int64
	needMinItemsCheck        bool
	minItems                 int64
	needUniqueItemsCheck     bool
}

func (av *arrayValidator) loadValidator() error {
	var err error

	if err = av.loadItems(); err != nil {
		return err
	}
	if err = av.loadAdditionalItem(); err != nil {
		return err
	}
	if err = av.loadMaxItems(); err != nil {
		return err
	}
	if err = av.loadMinItems(); err != nil {
		return err
	}
	if err = av.loadUniqueItems(); err != nil {
		return err
	}
	if err = av.loadContains(); err != nil {
		return err
	}
	return nil
}

func (av *arrayValidator) doValidate(jv *jvalue.V) (err error) {
	if jv.JType != jutil.JTArray {
		return nil
	}

	arr := jv.Value.([]*jvalue.V)
	if err = av.vItems(arr); err != nil {
		return err
	}

	if err = av.vMaxItems(arr); err != nil {
		return err
	}

	if err = av.vMinItems(arr); err != nil {
		return err
	}

	if err = av.vUniqueItems(arr); err != nil {
		return err
	}

	if err = av.vContains(arr); err != nil {
		return err
	}

	return nil
}

func (av *arrayValidator) loadItems() error {
	v, ok := av.schemaImpl.jvMap["items"]
	if !ok {
		av.needItemsCheck = false
		av.needItemsArrCheck = false
		return nil
	}

	if v.JType == jutil.JTArray {
		jvarr := v.Value.([]*jvalue.V)
		for i := 0; i < len(jvarr); i++ {
			key := arrayKey(av.schemaImpl.Schema, "items", i)
			s, err := parse(key, jvarr[i], av.schemaImpl.Schema,
				av.schemaImpl.Schema.root)
			if err != nil {
				return err
			}
			av.itemsArr = append(av.itemsArr, s)
		}
		av.needItemsArrCheck = true
		av.schemaImpl.svMap["items"] = newSchemaVArr(av.itemsArr)
		return nil
	}

	key := schemaKey(av.schemaImpl.Schema, "items")
	s, err := parse(key, v, av.schemaImpl.Schema,
		av.schemaImpl.Schema.root)
	if err != nil {
		return err
	}
	av.needItemsCheck = true
	av.items = s
	av.schemaImpl.svMap["items"] = newSchemaVSma(av.items)
	return nil
}

func (av *arrayValidator) loadAdditionalItem() error {
	v, ok := av.schemaImpl.jvMap["additionalItems"]
	if !ok {
		av.needAdditionalItemsCheck = false
		return nil
	}

	var err error
	key := schemaKey(av.schemaImpl.Schema, "additionalItems")
	av.additionalItems, err = parse(key, v, av.schemaImpl.Schema,
		av.schemaImpl.Schema.root)
	if err != nil {
		return err
	}

	av.schemaImpl.svMap["additionalItems"] = newSchemaVSma(av.additionalItems)
	av.needAdditionalItemsCheck = true
	return nil
}

func (av *arrayValidator) loadMaxItems() error {
	v, ok := av.schemaImpl.jvMap["maxItems"]
	if !ok {
		av.needMaxItemsCheck = false
		return nil
	}

	if v.JType != jutil.JTInteger {
		av.needMaxItemsCheck = false
		return fmt.Errorf("array maxItems must be integer")
	}

	av.maxItems, _ = v.GetInteger()
	if av.maxItems < 0 {
		return errors.New("'maxItems' is < 0")
	}
	av.needMaxItemsCheck = true
	return nil
}

func (av *arrayValidator) loadMinItems() error {
	v, ok := av.schemaImpl.jvMap["minItems"]
	if !ok {
		av.needMinItemsCheck = false
		return nil
	}

	if v.JType != jutil.JTInteger {
		av.needMinItemsCheck = false
		return fmt.Errorf("array minItems must be integer")
	}

	av.minItems, _ = v.GetInteger()
	if av.minItems < 0 {
		return errors.New("'minItems' is < 0")
	}
	av.needMinItemsCheck = true
	return nil
}

func (av *arrayValidator) loadUniqueItems() error {
	v, ok := av.schemaImpl.jvMap["uniqueItems"]
	if !ok {
		av.needUniqueItemsCheck = false
		return nil
	}

	if v.JType != jutil.JTBoolean {
		av.needUniqueItemsCheck = false
		return fmt.Errorf("array uniqueItems must be bool")
	}

	av.needUniqueItemsCheck, _ = v.GetBool()
	return nil
}

func (av *arrayValidator) loadContains() error {
	v, ok := av.schemaImpl.jvMap["contains"]
	if !ok {
		av.needContainsCheck = false
		return nil
	}

	key := schemaKey(av.schemaImpl.Schema, "contains")
	jv, err := parse(key, v, av.schemaImpl.Schema, av.schemaImpl.Schema.root)
	if err != nil {
		return err
	}

	av.contains = jv
	av.needContainsCheck = true
	av.schemaImpl.svMap["contains"] = newSchemaVSma(av.contains)
	return nil
}

func (av *arrayValidator) vItems(arr []*jvalue.V) (err error) {
	if !av.needItemsArrCheck && !av.needItemsCheck {
		return nil
	}

	if av.needItemsCheck {
		for i := 0; i < len(arr); i++ {
			if err = av.items.MatchJValue(arr[i]); err != nil {
				return err
			}
		}
		return nil
	}

	for i := 0; i < len(av.itemsArr) && i < len(arr); i++ {
		if err = av.itemsArr[i].MatchJValue(arr[i]); err != nil {
			return err
		}
	}

	// Checks for additional items.
	if len(arr) > len(av.itemsArr) && av.needAdditionalItemsCheck {
		for i := len(av.itemsArr); i < len(arr); i++ {
			if err = av.additionalItems.MatchJValue(arr[i]); err != nil {
				return err
			}
		}
	}
	return nil
}

func (av *arrayValidator) vMaxItems(arr []*jvalue.V) (err error) {
	if !av.needMaxItemsCheck {
		return nil
	}

	if len(arr) > int(av.maxItems) {
		return fmt.Errorf("failed to match maxItems. schema is %s",
			av.schemaImpl.Schema.key)
	}

	return nil
}

func (av *arrayValidator) vMinItems(arr []*jvalue.V) (err error) {
	if !av.needMinItemsCheck {
		return nil
	}

	if len(arr) < int(av.minItems) {
		return fmt.Errorf("failed to match minItems. schema is %s",
			av.schemaImpl.Schema.key)
	}

	return nil
}

func (av *arrayValidator) vUniqueItems(arr []*jvalue.V) (err error) {
	if !av.needUniqueItemsCheck {
		return nil
	}

	for i := 0; i < len(arr); i++ {
		for j := i + 1; j < len(arr); j++ {
			if err = arr[i].IsEqual(arr[j]); err == nil {
				return fmt.Errorf("uniqueItems failed. i is %d. "+
					"j is %d. schema is %s",
					i, j, av.schemaImpl.Schema.key)
			}
		}
	}

	return nil
}

func (av *arrayValidator) vContains(arr []*jvalue.V) (err error) {
	if !av.needContainsCheck {
		return nil
	}

	for i := 0; i < len(arr); i++ {
		if err = av.contains.MatchJValue(arr[i]); err == nil {
			return nil
		}
	}

	return fmt.Errorf("contains failed. schema is %s", av.schemaImpl.Schema.key)
}
