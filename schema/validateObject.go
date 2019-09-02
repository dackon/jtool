package schema

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type patternProperty struct {
	r *regexp.Regexp
	s *Schema
}

type objectValidator struct {
	schemaImpl *schemaImpl

	// For objects.
	needMaxPropertiesCheck        bool
	maxProperties                 int64
	needMinPropertiesCheck        bool
	minProperties                 int64
	needRequiredCheck             bool
	required                      []string
	needPropertiesCheck           bool
	properties                    map[string]*Schema
	needPatternPropertiesCheck    bool
	patternProperites             []*patternProperty
	needAdditionalPropertiesCheck bool
	additionalProperties          *Schema
	needDependenciesCheck         bool
	dependencies                  map[string]*dependency
	needPropertyNamesCheck        bool
	propertyNames                 *Schema
}

func (ov *objectValidator) loadValidator() error {
	var err error
	if err = ov.loadAdditionalProperties(); err != nil {
		return err
	}

	if err = ov.loadDependenciesProperties(); err != nil {
		return err
	}

	if err = ov.loadMaxProperties(); err != nil {
		return err
	}

	if err = ov.loadMinProperties(); err != nil {
		return err
	}

	if err = ov.loadPatternProperties(); err != nil {
		return err
	}

	if err = ov.loadProperties(); err != nil {
		return err
	}

	if err = ov.loadPropertyNames(); err != nil {
		return err
	}

	if err = ov.loadRequired(); err != nil {
		return err
	}

	return nil
}

func (ov *objectValidator) doValidate(jv *jvalue.V) (err error) {
	if jv.JType != jutil.JTObject {
		return nil
	}
	mv := jv.Value.(map[string]*jvalue.V)

	// additionalProperties checking is included in properties and patternPropertes.

	if err = ov.vDependencies(jv); err != nil {
		return err
	}

	if err = ov.vMaxProperties(mv); err != nil {
		return err
	}

	if err = ov.vMinProperties(mv); err != nil {
		return err
	}

	if err = ov.vProperties(mv); err != nil {
		return err
	}

	if err = ov.vPropertyNames(mv); err != nil {
		return err
	}

	if err = ov.vRequired(mv); err != nil {
		return err
	}

	return nil
}

func (ov *objectValidator) loadMaxProperties() error {
	v, ok := ov.schemaImpl.jvMap["maxProperties"]
	if !ok {
		ov.needMaxPropertiesCheck = false
		return nil
	}

	if v.JType != jutil.JTInteger {
		ov.needMaxPropertiesCheck = false
		return fmt.Errorf("maxProperties must be integer")
	}

	ov.maxProperties, _ = v.GetInteger()
	if ov.maxProperties < 0 {
		return errors.New("maxProperties value must be >= 0")
	}
	ov.needMaxPropertiesCheck = true
	return nil
}

func (ov *objectValidator) loadMinProperties() error {
	v, ok := ov.schemaImpl.jvMap["minProperties"]
	if !ok {
		ov.needMinPropertiesCheck = false
		return nil
	}

	if v.JType != jutil.JTInteger {
		ov.needMinPropertiesCheck = false
		return fmt.Errorf("minProperties must be integer")
	}

	ov.minProperties, _ = v.GetInteger()
	if ov.minProperties < 0 {
		return errors.New("minProperties value must be >= 0")
	}
	ov.needMinPropertiesCheck = true
	return nil
}

func (ov *objectValidator) loadRequired() error {
	v, ok := ov.schemaImpl.jvMap["required"]
	if !ok {
		ov.needRequiredCheck = false
		return nil
	}

	var err error
	ov.required, err = v.GetStringArr()
	if err != nil {
		ov.needRequiredCheck = false
		return fmt.Errorf("required must be string array")
	}

	if len(ov.required) == 0 {
		ov.needRequiredCheck = false
		return nil
	}

	for i := 0; i < len(ov.required); i++ {
		for j := i + 1; j < len(ov.required); j++ {
			if ov.required[i] == ov.required[j] {
				return errors.New("'required' must be unique")
			}
		}
	}

	ov.needRequiredCheck = true
	return nil
}

func (ov *objectValidator) loadProperties() error {
	v, ok := ov.schemaImpl.jvMap["properties"]
	if !ok {
		ov.needPropertiesCheck = false
		return nil
	}

	if v.JType != jutil.JTObject {
		ov.needPropertiesCheck = false
		return fmt.Errorf("properties must be object")
	}

	rv := v.Value.(map[string]*jvalue.V)
	ov.properties = make(map[string]*Schema, len(rv))
	for k, v := range rv {
		key := objectKey(ov.schemaImpl.Schema, "properties", k)
		s, err := parse(key, v, ov.schemaImpl.Schema,
			ov.schemaImpl.Schema.root)
		if err != nil {
			return err
		}
		ov.properties[k] = s
	}

	ov.schemaImpl.svMap["properties"] = newSchemaVObj(ov.properties)
	ov.needPropertiesCheck = true
	return nil
}

func (ov *objectValidator) loadPatternProperties() error {
	v, ok := ov.schemaImpl.jvMap["patternProperties"]
	if !ok {
		ov.needPatternPropertiesCheck = false
		return nil
	}

	if v.JType != jutil.JTObject {
		ov.needPatternPropertiesCheck = false
		return fmt.Errorf("patternProperties must be object")
	}

	rv := v.Value.(map[string]*jvalue.V)
	ov.patternProperites = make([]*patternProperty, 0, len(rv))
	ksMap := make(map[string]*Schema, len(rv))
	for k, v := range rv {
		key := objectKey(ov.schemaImpl.Schema, "patternProperties", k)
		s, err := parse(key, v, ov.schemaImpl.Schema,
			ov.schemaImpl.Schema.root)
		if err != nil {
			return err
		}

		r, err := regexp.Compile(k)
		if err != nil {
			return err
		}

		pp := &patternProperty{
			r: r,
			s: s,
		}
		ov.patternProperites = append(ov.patternProperites, pp)
		ksMap[k] = s
	}

	ov.schemaImpl.svMap["patternProperites"] = newSchemaVObj(ksMap)
	ov.needPatternPropertiesCheck = true
	return nil
}

func (ov *objectValidator) loadAdditionalProperties() error {
	v, ok := ov.schemaImpl.jvMap["additionalProperties"]
	if !ok {
		ov.needAdditionalPropertiesCheck = false
		return nil
	}

	key := schemaKey(ov.schemaImpl.Schema, "additionalProperties")
	s, err := parse(key, v, ov.schemaImpl.Schema, ov.schemaImpl.Schema.root)
	if err != nil {
		return err
	}

	ov.needAdditionalPropertiesCheck = true
	ov.additionalProperties = s
	return nil
}

func (ov *objectValidator) loadDependenciesProperties() error {
	v, ok := ov.schemaImpl.jvMap["dependencies"]
	if !ok {
		ov.needDependenciesCheck = false
		return nil
	}

	if v.JType != jutil.JTObject {
		ov.needDependenciesCheck = false
		return fmt.Errorf("dependencies must be object")
	}

	var err error

	mv := v.Value.(map[string]*jvalue.V)
	ov.dependencies = make(map[string]*dependency, len(mv))
	dmap := make(map[string]*Schema, len(mv))

	for k, d := range mv {
		dp := &dependency{}
		if d.JType == jutil.JTArray {
			dp.keys, err = d.GetStringArr()
			if err != nil {
				return err
			}
		} else {
			key := objectKey(ov.schemaImpl.Schema, "dependencies", k)
			dp.s, err = parse(key, d, ov.schemaImpl.Schema,
				ov.schemaImpl.Schema.root)
			if err != nil {
				return err
			}

			dmap[k] = dp.s
		}
		ov.dependencies[k] = dp
	}

	if len(dmap) > 0 {
		ov.schemaImpl.svMap["dependencies"] = newSchemaVObj(dmap)
	}
	ov.needDependenciesCheck = true
	return nil
}

func (ov *objectValidator) loadPropertyNames() error {
	v, ok := ov.schemaImpl.jvMap["propertyNames"]
	if !ok {
		ov.needPropertyNamesCheck = false
		return nil
	}

	var err error
	key := schemaKey(ov.schemaImpl.Schema, "propertyNames")
	ov.propertyNames, err = parse(key, v, ov.schemaImpl.Schema,
		ov.schemaImpl.Schema.root)
	if err != nil {
		return err
	}

	ov.schemaImpl.svMap["propertyNames"] = newSchemaVSma(ov.propertyNames)
	ov.needPropertyNamesCheck = true
	return nil
}

func (ov *objectValidator) vMaxProperties(mjv map[string]*jvalue.V) (
	err error) {
	if !ov.needMaxPropertiesCheck {
		return nil
	}

	if len(mjv) > int(ov.maxProperties) {
		return fmt.Errorf("maxProperties failed. schema is %s",
			ov.schemaImpl.Schema.key)
	}

	return nil
}

func (ov *objectValidator) vMinProperties(mjv map[string]*jvalue.V) (
	err error) {
	if !ov.needMinPropertiesCheck {
		return nil
	}

	if len(mjv) < int(ov.minProperties) {
		return fmt.Errorf("minProperties failed. schema is %s",
			ov.schemaImpl.Schema.key)
	}

	return nil
}

func (ov *objectValidator) vRequired(mjv map[string]*jvalue.V) (
	err error) {

	if !ov.needRequiredCheck {
		return nil
	}

	for i := 0; i < len(ov.required); i++ {
		_, ok := mjv[ov.required[i]]
		if !ok {
			return fmt.Errorf("required failed. key is %s. schema is %s",
				ov.required[i], ov.schemaImpl.Schema.key)
		}
	}

	return nil
}

func (ov *objectValidator) vProperties(mjv map[string]*jvalue.V) (
	err error) {

	var s *Schema
	var ok bool
	for k, v := range mjv {

		// ok means whether the k is found in properties or patternProperties.
		ok = false

		// If need properties check, get the schema corresponding to the key.
		if ov.needPropertiesCheck {
			s, ok = ov.properties[k]
			if ok {
				if err = s.MatchJValue(v); err != nil {
					return err
				}
			}
		}

		// Even if v is valid in properties check, we need to check
		// patternProperties again.
		if ov.needPatternPropertiesCheck {
			for i := 0; i < len(ov.patternProperites); i++ {
				p := ov.patternProperites[i]
				if p.r.MatchString(k) {
					if err = p.s.MatchJValue(v); err != nil {
						return err
					}

					ok = true
				}
			}
		}

		// If not found the key in properties and patternProperties, do
		// addtionalProperties check.
		if !ok && ov.needAdditionalPropertiesCheck {
			if err = ov.additionalProperties.MatchJValue(v); err != nil {
				return err
			}
		}
	}

	return nil
}

func (ov *objectValidator) vDependencies(jv *jvalue.V) (
	err error) {

	if !ov.needDependenciesCheck {
		return nil
	}

	mjv := jv.Value.(map[string]*jvalue.V)

	for k := range mjv {
		d, ok := ov.dependencies[k]
		if !ok {
			continue
		}

		if d.s == nil {
			// Use key array.
			for i := 0; i < len(d.keys); i++ {
				_, exist := mjv[d.keys[i]]
				if !exist {
					return fmt.Errorf("dependencies failed. key is %s. "+
						"schema is %s", d.keys[i], ov.schemaImpl.Schema.key)
				}
			}
		} else {
			// Use schema.
			if err = d.s.MatchJValue(jv); err != nil {
				return err
			}
		}
	}

	return nil
}

func (ov *objectValidator) vPropertyNames(mjv map[string]*jvalue.V) (
	err error) {

	if !ov.needPropertyNamesCheck {
		return nil
	}

	for k := range mjv {
		kstr := fmt.Sprintf("\"%s\"", k)
		if err = ov.propertyNames.Match([]byte(kstr)); err != nil {
			return err
		}
	}

	return nil
}
