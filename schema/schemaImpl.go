package schema

import (
	"fmt"
	"net/url"
	"regexp"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type dependency struct {
	keys []string
	s    *Schema
}

type valueType struct {
	Type    jutil.JType
	TypeArr []jutil.JType
}

type schemaImpl struct {
	Schema      *Schema
	jvMap       map[string]*jvalue.V
	svMap       map[string]*schemaValue
	emptySchema bool

	hasType bool
	Type    valueType

	id      *url.URL
	baseURI *url.URL

	definitions *definitions
	validators  []validatorInf
}

var (
	gFragmentReg *regexp.Regexp
)

func init() {
	gFragmentReg = regexp.MustCompile(`^#[a-zA-Z]+[0-9a-zA-Z\-\_\:\.]*$`)
}

func newSchemaImpl(s *Schema, jvMap map[string]*jvalue.V) *schemaImpl {
	return &schemaImpl{
		jvMap:  jvMap,
		Schema: s,
		svMap:  make(map[string]*schemaValue, 8),
	}
}

func (si *schemaImpl) parse() error {
	if len(si.jvMap) == 0 {
		si.emptySchema = true
		return nil
	}

	var err error
	err = si.setType()
	if err != nil {
		return err
	}

	if err = si.setID(); err != nil {
		return err
	}

	// Load definitions.
	si.definitions = newDefinitions(si)
	if err = si.definitions.loadDefinitions(); err != nil {
		return err
	}

	enumv := &enumValidator{schemaImpl: si}
	si.validators = append(si.validators, enumv)

	constv := &constValidator{schemaImpl: si}
	si.validators = append(si.validators, constv)

	numv := &numberValidator{schemaImpl: si}
	si.validators = append(si.validators, numv)

	intv := &integerValidator{schemaImpl: si}
	si.validators = append(si.validators, intv)

	strv := &stringValidator{schemaImpl: si}
	si.validators = append(si.validators, strv)

	arrv := &arrayValidator{schemaImpl: si}
	si.validators = append(si.validators, arrv)

	objv := &objectValidator{schemaImpl: si}
	si.validators = append(si.validators, objv)

	cdlv := &condValidator{schemaImpl: si}
	si.validators = append(si.validators, cdlv)

	bllv := &boolLogicValidator{schemaImpl: si}
	si.validators = append(si.validators, bllv)

	// TODO: implement formatValidator.

	for i := 0; i < len(si.validators); i++ {
		v := si.validators[i]
		if err = v.loadValidator(); err != nil {
			return err
		}
	}
	return nil
}

func (si *schemaImpl) setType() (err error) {
	v, ok := si.jvMap["type"]
	if !ok {
		si.hasType = false
		return nil
	}

	stype := &si.Type
	if v.JType == jutil.JTArray {
		sarr, err := v.GetStringArr()
		if err != nil {
			return fmt.Errorf("schema %s type is not string array",
				si.Schema.key)
		}
		for _, s := range sarr {
			t, err := jutil.IsValidJSONType(s)
			if err != nil {
				return fmt.Errorf("bad schema type. schema is %s",
					si.Schema.key)
			}
			stype.TypeArr = append(stype.TypeArr, t)
		}
	} else if v.JType == jutil.JTString {
		str, _ := v.GetString()
		t, err := jutil.IsValidJSONType(str)
		if err != nil {
			return fmt.Errorf("bad schema type. schema is %s",
				si.Schema.key)
		}
		stype.Type = t
	} else {
		return fmt.Errorf("bad schema type. schema is %s", si.Schema.key)
	}

	si.hasType = true
	return nil
}

func (si *schemaImpl) setID() (err error) {
	id, ok := si.jvMap["$id"]
	if !ok {
		if si.Schema.parent == nil {
			// This is root schema.
			si.id, err = url.Parse("/")
			if err != nil {
				return fmt.Errorf("url.Parse failed. Err is %s. schema is %s",
					err, si.Schema.key)
			}
			si.baseURI = si.id
			// ref RFC 6901, section 6
			si.Schema.schemaMap["#"] = si.Schema
		}
		return
	}

	if id.JType != jutil.JTString {
		return fmt.Errorf("the value of $id must be string. schema is %s",
			si.Schema.key)
	}

	sv, _ := id.GetString()
	si.id, err = url.Parse(sv)
	if err != nil {
		return fmt.Errorf("url.Parse failed. Err is %s. schema is %s",
			err, si.Schema.key)
	}

	p := si.Schema.parent
	if p == nil {
		// s.Schema is root Schema.
		si.baseURI = si.id
		empty, _ := url.Parse("")
		si.baseURI = si.baseURI.ResolveReference(empty)
		si.Schema.schemaMap[si.baseURI.String()] = si.Schema
		return nil
	}

	if p.schemaImpl.id == nil {
		si.id = nil
		si.baseURI = nil
		return nil
	}

	pid := p.schemaImpl.id.String()
	if len(pid) > 0 && pid[0] == '#' {
		// If parent is fragment, do not need to resolve the URI.
		si.id = nil
		si.baseURI = nil
		return nil
	}

	sid := si.id.String()
	if sid == "" || sid == "#" {
		return fmt.Errorf("$id of subschema must not be \"\" or \"#\". "+
			"schema is %s", si.Schema.key)
	}

	var absURI *url.URL
	if sid[0] == '#' {
		// The $id is fragment.
		if !gFragmentReg.MatchString(sid) {
			return fmt.Errorf("bad $id format, $id must match regexp "+
				"\"^#[a-zA-Z]+[0-9a-zA-Z-_:.]*$\". schema is %s", si.Schema.key)
		}

		si.baseURI = p.schemaImpl.baseURI
		absURI = si.baseURI.ResolveReference(si.id)
	} else {
		si.baseURI = p.schemaImpl.baseURI.ResolveReference(si.id)
		absURI = si.baseURI
	}

	_, ok = si.Schema.root.schemaMap[absURI.String()]
	if ok {
		return fmt.Errorf("duplicated ID %s. schema is %s", sid, si.Schema.key)
	}

	si.Schema.root.schemaMap[absURI.String()] = si.Schema
	return nil
}

func (si *schemaImpl) validType(jv *jvalue.V) (err error) {
	if !si.hasType {
		return nil
	}

	stype := &si.Type
	if stype.Type != jutil.JTNone {
		if stype.Type != jv.JType {
			if stype.Type == jutil.JTNumber && jv.JType == jutil.JTInteger {
				return nil
			}
			return fmt.Errorf("match type failed. schema is %s", si.Schema.key)
		}
		return nil
	}

	for i := 0; i < len(stype.TypeArr); i++ {
		if stype.TypeArr[i] == jv.JType {
			return nil
		}
		if stype.TypeArr[i] == jutil.JTNumber && jv.JType == jutil.JTInteger {
			return nil
		}
	}

	return fmt.Errorf("match type failed. schema is %s", si.Schema.key)
}

func (si *schemaImpl) isType(t jutil.JType) bool {
	if !si.hasType {
		return false
	}

	stype := &si.Type
	if stype.Type != jutil.JTNone {
		if stype.Type != t {
			return false
		}
		return true
	}

	for i := 0; i < len(stype.TypeArr); i++ {
		if stype.TypeArr[i] == t {
			return true
		}
	}

	return false
}

func (si *schemaImpl) match(jv *jvalue.V) (err error) {
	if jv == nil {
		return fmt.Errorf("match failed. Target is nil")
	}

	// If schema doesn't have a type, jv is valid. If schema has a type,
	// it must match the type of jv.
	err = si.validType(jv)
	if err != nil {
		return err
	}

	// Here, schema may not have a type, or the type of jv must match the
	// type of schema.

	for i := 0; i < len(si.validators); i++ {
		if err = si.validators[i].doValidate(jv); err != nil {
			return err
		}
	}

	return nil
}
