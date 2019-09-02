package schema

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

// Schema ...
type Schema struct {
	kind       schemaKind
	boolValue  bool
	schemaImpl *schemaImpl
	parent     *Schema
	root       *Schema
	isRef      bool
	ref        string
	refSchema  *Schema
	key        string

	// <Absolute URI, *Schema>
	schemaMap map[string]*Schema
}

// Parse ...
func Parse(raw json.RawMessage) (*Schema, error) {
	if len(raw) == 0 {
		return nil, ErrBadParam
	}

	jv, err := jvalue.ParseJSON(raw)
	if err != nil {
		return nil, err
	}

	return ParseJV(jv)
}

// ParseJV ...
func ParseJV(jv *jvalue.V) (*Schema, error) {
	if jv == nil {
		return nil, fmt.Errorf("Bad parameter")
	}

	s, err := parse("/", jv, nil, nil)
	if err != nil {
		return nil, err
	}

	if s.kind == schemaBool {
		return s, nil
	}

	if err = s.resolveRef(nil); err != nil {
		return nil, err
	}

	return s, nil
}

func parse(key string, jv *jvalue.V, parent *Schema, root *Schema) (
	*Schema, error) {
	schema := &Schema{
		parent: parent,
		root:   root,
		key:    key,
	}

	if root == nil {
		schema.root = schema
	}

	switch jv.JType {
	case jutil.JTObject:

		// Schema kind is schemaObject.
		schema.kind = schemaObject
		mjv := jv.Value.(map[string]*jvalue.V)

		v, ok := mjv["$ref"]
		if ok {

			// $ref schema.

			if parent == nil {
				return nil, fmt.Errorf("root schema cannot be $ref")
			}

			if v.JType != jutil.JTString {
				return nil, fmt.Errorf("'$ref' must be string")
			}

			// Update ref info.
			schema.isRef = true
			schema.ref, _ = v.GetString()
			return schema, nil
		}

		// If here, normal schemaObject.

		schema.schemaImpl = newSchemaImpl(schema, mjv)
		if schema.root == schema {
			// This is root schema.
			schema.schemaMap = make(map[string]*Schema, 8)
		}

		if err := schema.schemaImpl.parse(); err != nil {
			return nil, err
		}
		return schema, nil
	case jutil.JTBoolean:

		// Schema kind is schemaBool.
		schema.kind = schemaBool
		v, _ := jv.GetBool()
		if v {
			schema.boolValue = true
		} else {
			schema.boolValue = false
		}

		return schema, nil
	}

	return nil, fmt.Errorf("failed to parse schema")
}

// resolveRef resolves $ref in schema. revArr is used to detect circle $ref.
func (s *Schema) resolveRef(revArr []*Schema) error {
	if s.kind == schemaBool {
		return nil
	}
	if !s.isRef {
		for _, v := range s.schemaImpl.svMap {
			switch v.svType {
			case svObject:
				rv := v.value.(map[string]*Schema)
				for _, y := range rv {
					if err := y.resolveRef(revArr); err != nil {
						return err
					}
				}
			case svArray:
				rv := v.value.([]*Schema)
				for _, y := range rv {
					if err := y.resolveRef(revArr); err != nil {
						return err
					}
				}
			case svSchema:
				rv := v.value.(*Schema)
				if err := rv.resolveRef(revArr); err != nil {
					return err
				}
			}
		}
		return nil
	}

	// Here, the schema is ref schema.

	if s.refSchema != nil {
		return nil
	}

	ts, err := s.getSchema()
	if err != nil {
		return err
	}

	if !ts.isRef || (ts.isRef && ts.refSchema != nil) {
		s.refSchema = ts
		return nil
	}

	// Add self to detect circle reference.
	revArr = append(revArr, s)
	for i := 0; i < len(revArr); i++ {
		if revArr[i] == ts {
			return fmt.Errorf("'$ref' circle detected: [%s, %s]",
				revArr[i].ref, ts.ref)
		}
	}

	err = ts.resolveRef(revArr)
	if err != nil {
		return err
	}

	if ts.refSchema == nil {
		return fmt.Errorf("failed to find $ref %s", s.ref)
	}
	s.refSchema = ts
	return nil
}

// Match ...
func (s *Schema) Match(raw json.RawMessage) error {
	switch s.kind {
	case schemaBool:
		if s.boolValue == true {
			return nil
		}
		return fmt.Errorf("failed to match schema %s", s.key)
	case schemaObject:
		if s.isRef {
			if err := s.refSchema.Match(raw); err != nil {
				return fmt.Errorf("failed to match schema %s", s.key)
			}
			return nil
		}

		jv, err := jvalue.ParseJSON(raw)
		if err != nil {
			return err
		}
		return s.schemaImpl.match(jv)
	}

	panic("Match must not be here")
}

// MatchJValue ...
func (s *Schema) MatchJValue(jv *jvalue.V) error {
	switch s.kind {
	case schemaBool:
		if s.boolValue == true {
			return nil
		}
		return fmt.Errorf("failed to match schema %s", s.key)
	case schemaObject:
		if s.isRef {
			if err := s.refSchema.MatchJValue(jv); err != nil {
				return fmt.Errorf("failed to match schema %s", s.key)
			}
			return nil
		}
		return s.schemaImpl.match(jv)
	}

	panic("Match must not be here")
}

func (s *Schema) isValidJSONValue(jv *jvalue.V) error {
	if s.kind == schemaBool {
		if s.boolValue == true {
			return nil
		}

		return fmt.Errorf("not valid json value. schema is %s", s.key)
	}

	return s.schemaImpl.match(jv)
}

func (s *Schema) getSchema() (*Schema, error) {
	refURI, err := url.Parse(s.ref)
	if err != nil {
		return nil, err
	}

	escapedRef := refURI.String()
	if escapedRef == "" || escapedRef == "#" {
		return s.root, nil
	}

	if len(escapedRef) > 1 && escapedRef[0] == '#' && escapedRef[1] == '/' {
		schema := findParentWithBaseURL(s)
		return findSchemaByJSONPointer(escapedRef, schema)
	}

	if escapedRef[0] == '#' {
		schema := findParentWithBaseURL(s)
		return findSchemaByFragment(refURI, schema)
	}

	// If here, the $ref has path, and may have fragment.

	// 1. Firstly, find the parent with the base URL
	schema := findParentWithBaseURL(s)
	absURI := schema.schemaImpl.baseURI.ResolveReference(refURI).String()
	ts, ok := s.root.schemaMap[absURI]
	if ok {
		return ts, nil
	}

	// 2. Not found by the URI, if the refURI has fragment and is json pointer,
	//    find the base schema firstly, then, find the target schema according
	//    to the json pointer.

	refURICopy := *refURI
	refURICopy.Fragment = ""

	if len(refURI.Fragment) > 0 && refURI.Fragment[0] != '/' {
		// Here, the fragment is json pointer.
		absURI = schema.schemaImpl.baseURI.ResolveReference(
			&refURICopy).String()
		ts, ok = s.root.schemaMap[absURI]
		if ok {
			ts, err = findSchemaByJSONPointer(refURI.Fragment, ts)
			if err == nil {
				return ts, nil
			}
		}
	}

	// 3. Not found in parent schema. Loop the schemaMap.

	for _, v := range s.root.schemaMap {
		absURI = v.schemaImpl.baseURI.ResolveReference(refURI).String()
		ts, ok = s.root.schemaMap[absURI]
		if ok {
			return ts, nil
		}
	}

	// 4. If here, and if the $ref has fragment and is not JSON pointer, then,
	//    we failed to find the schema.
	if len(refURI.Fragment) > 0 && refURI.Fragment[0] != '/' {
		return nil, fmt.Errorf("failed to find $ref %s", s.ref)
	}

	// 5. If here, we get the path firstly, and then, find the schema by that
	//    path. If we find the schema, we can find the referenced schema by
	//    JSON pointer against that schema.
	for _, v := range s.root.schemaMap {
		absURI = v.schemaImpl.baseURI.ResolveReference(&refURICopy).String()
		ts, ok = s.root.schemaMap[absURI]
		if ok {
			break
		}
	}

	if !ok {
		return nil, fmt.Errorf("failed to find $ref %s", s.ref)
	}

	return findSchemaByJSONPointer(refURI.Fragment, ts)
}

func findParentWithBaseURL(s *Schema) *Schema {
	schema := s.parent
	for {
		if schema == nil {
			return s.root
		}

		if schema.schemaImpl.baseURI == nil {
			schema = schema.parent
		} else {
			break
		}

		if schema.parent == nil {
			// schema is root
			break
		}
	}

	return schema
}

func findSchemaByJSONPointer(jp string, start *Schema) (*Schema, error) {
	elems, err := decodeJSONPointer(jp)
	if err != nil {
		return nil, err
	}

	if len(elems) == 0 {
		// If $ref is '#' or empty string, the referenced schema is root
		// schema.
		return start.root, nil
	}

	schema := start
	for i := 0; i < len(elems); {
		v, ok := schema.schemaImpl.svMap[elems[i]]
		if !ok {
			return nil, fmt.Errorf("failed to find key %s", elems[i])
		}

		if v.svType == svObject || v.svType == svArray {
			i++
		}

		schema, err = v.getSchema(elems[i])
		if err != nil {
			return nil, err
		}
		i++
	}

	return schema, nil
}

func findSchemaByFragment(fragment *url.URL, start *Schema) (*Schema, error) {
	absURI := start.schemaImpl.baseURI.ResolveReference(fragment).String()
	ts, ok := start.root.schemaMap[absURI]
	if !ok {
		return nil, fmt.Errorf("faild to find $ref %s", fragment.String())
	}
	return ts, nil
}
