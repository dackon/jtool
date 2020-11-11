package schema8

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

func doParse(raw json.RawMessage, canonicalURL string) (
	*Schema, error) {
	v, err := jvalue.ParseJSON(raw)
	if err != nil {
		return nil, err
	}

	defaultBaseURI := canonicalURL
	if defaultBaseURI == "" {
		defaultBaseURI = "https://default.uri"
	}

	bu, err := url.Parse(defaultBaseURI)
	if err != nil || !bu.IsAbs() {
		return nil, fmt.Errorf("canonicalURL is not a valid absolute URL")
	}

	// Remove fragment
	bu.Fragment = ""
	defaultBaseURI = bu.String()

	schema := &Schema{schemaJar: NewJar()}

	switch v.JType {
	case jutil.JTBoolean:
		schema.root = &schemaNode{
			snType:        sntBoolSchema,
			jtType:        jutil.JTBoolean,
			value:         v.Value.(bool),
			baseURI:       defaultBaseURI,
			canonicalURIs: []string{defaultBaseURI + "#"},
		}
	case jutil.JTObject:
		mv := v.Value.(map[string]*jvalue.V)
		id, ok := mv["$id"]
		if !ok {
			mv["$id"] = &jvalue.V{
				JType:       jutil.JTString,
				Value:       defaultBaseURI,
				Parent:      v,
				KeyInParent: "$id",
			}
		} else {
			s, err := id.GetString()
			if err != nil {
				return nil, fmt.Errorf("$id must be a string")
			}

			u, err := url.Parse(s)
			if err != nil {
				return nil, fmt.Errorf("failed to parse base URI")
			}

			if u.Scheme == "" {
				u = bu.ResolveReference(u)
				mv["$id"] = &jvalue.V{
					JType:       jutil.JTString,
					Value:       u.String(),
					Parent:      v,
					KeyInParent: "$id",
				}
			}
		}

		head := &schemaNode{
			schema: schema,
		}

		pn := newPathNodeString("")
		sn, err := doParseSchema(v, head, pn)
		if err != nil {
			return nil, err
		}

		err = sn.finalize()
		schema.root = sn

		err = schema.schemaJar.LinkRef()
		if err != nil {
			return nil, err
		}

	default:
		return nil, ErrNotSchema
	}

	return schema, nil
}

func doParseSchema(v *jvalue.V, parent *schemaNode, pn *pathNode) (
	*schemaNode, error) {

	sn := &schemaNode{
		jtType: v.JType,
		path:   copyPath(parent, pn),
		value:  v.Value,
		schema: parent.schema,
		kvMap:  make(map[string]*schemaNode, 8),
	}

	if v.Parent != nil {
		sn.parent = parent
	}

	if v == nil {
		return nil, errWithPath(ErrNotSchema, sn)
	}

	if v.JType != jutil.JTBoolean && v.JType != jutil.JTObject {
		return nil, errWithPath(ErrNotSchema, sn)
	}

	switch v.JType {
	case jutil.JTBoolean:
		sn.snType = sntBoolSchema

	case jutil.JTObject:
		sn.snType = sntObjectSchema
		mv := v.Value.(map[string]*jvalue.V)
		if err := doParseSchemaFields(sn, mv); err != nil {
			return nil, err
		}
	}

	return sn, nil
}

func doParseObjectSchema(v *jvalue.V,
	parent *schemaNode, pn *pathNode) (*schemaNode, error) {
	sn := &schemaNode{
		snType: sntNotSchema,
		jtType: v.JType,
		schema: parent.schema,
		parent: parent,
		path:   copyPath(parent, pn),
		value:  v.Value,
		kvMap:  make(map[string]*schemaNode, 8),
	}

	if v.JType != jutil.JTObject {
		return nil, errWithPath(fmt.Errorf("value is not object"), sn)
	}

	mv := v.Value.(map[string]*jvalue.V)
	for k, v := range mv {
		pn := newPathNodeString(k)
		node, err := doParseSchema(v, sn, pn)
		if err != nil {
			return nil, err
		}
		sn.kvMap[k] = node
	}

	return sn, nil
}

func doParseArraySchema(v *jvalue.V,
	parent *schemaNode, pn *pathNode) (*schemaNode, error) {
	sn := &schemaNode{
		snType: sntNotSchema,
		jtType: v.JType,
		schema: parent.schema,
		parent: parent,
		path:   copyPath(parent, pn),
		value:  v.Value,
	}

	if v.JType != jutil.JTArray {
		return nil, errWithPath(fmt.Errorf("value is not array"), sn)
	}

	arr := v.Value.([]*jvalue.V)
	for i, v := range arr {
		pn := newPathNodeInt(i)
		node, err := doParseSchema(v, sn, pn)
		if err != nil {
			return nil, err
		}
		sn.nodeArr = append(sn.nodeArr, node)
	}

	return sn, nil
}

func doParseJSON(v *jvalue.V,
	parent *schemaNode, pn *pathNode) (*schemaNode, error) {
	sn := &schemaNode{
		jtType: v.JType,
		snType: sntNotSchema,
		schema: parent.schema,
		parent: parent,
		path:   copyPath(parent, pn),
		value:  v.Value,
		kvMap:  make(map[string]*schemaNode, 8),
	}

	switch v.JType {
	case jutil.JTObject:
		mv := v.Value.(map[string]*jvalue.V)
		for k, v := range mv {
			pn := newPathNodeString(k)
			node, err := doParseJSON(v, sn, pn)
			if err != nil {
				return nil, err
			}
			sn.kvMap[k] = node
		}
	case jutil.JTArray:
		arr := v.Value.([]*jvalue.V)
		for i, v := range arr {
			pn := newPathNodeInt(i)
			node, err := doParseJSON(v, sn, pn)
			if err != nil {
				return nil, err
			}
			sn.nodeArr = append(sn.nodeArr, node)
		}
	}

	return sn, nil
}

func doParseSchemaFields(sn *schemaNode, mv map[string]*jvalue.V) error {
	var err error
	for key, value := range mv {
		if err = loadCore(key, value, sn); err != nil {
			return err
		}

		if err = loadAssersion(key, value, sn); err != nil {
			return err
		}

		if err = loadApplicator(key, value, sn); err != nil {
			return err
		}
	}

	finalizeIfThenElse(sn)
	finalizeContains(sn)
	finalizeAdditionalItems(sn)
	finalizeAdditionalProperties(sn)
	return nil
}

func loadCore(key string, value *jvalue.V, sn *schemaNode) error {
	switch key {
	// core
	case "$schema":
		s, err := value.GetString()
		if err != nil {
			return errWithPath(fmt.Errorf("$schema must be a string"), sn)
		}

		if sn.parent != nil {
			return errWithPath(fmt.Errorf(
				"$schema must not be in subschemas"), sn)
		}

		u, err := url.Parse(s)
		if err != nil || u.Scheme == "" {
			return errWithPath(fmt.Errorf(
				"$schema must be a valid absolute url"), sn)
		}

		sn.metaSchema = s
	case "$id":
		s, err := value.GetString()
		if err != nil {
			return errWithPath(fmt.Errorf("$id must be a string"), sn)
		}

		u, err := url.Parse(s)
		if err != nil {
			return errWithPath(fmt.Errorf(
				"$id must be a valid url-reference"), sn)
		}

		if u.Fragment != "" {
			return errWithPath(fmt.Errorf(
				"$id '%s' must not have none-empty fragment", u.String()), sn)
		}

		// Normalize $id to avoid any surprises.
		sn.id = u.String()
		sn.idURL = u

	case "$anchor":
		s, err := value.GetString()
		if err != nil {
			return errWithPath(fmt.Errorf("$anchor must be a string"), sn)
		}
		if !gAnchorReg.MatchString(s) {
			return errWithPath(fmt.Errorf("bad $anchor %s", s), sn)
		}
		sn.anchor = url.PathEscape(s)

	case "$ref":
		s, err := value.GetString()
		if err != nil {
			return errWithPath(fmt.Errorf("$ref must be a string"), sn)
		}

		u, err := url.Parse(s)
		if err != nil {
			return errWithPath(
				fmt.Errorf("$ref must be a valid url-reference"), sn)
		}

		sn.ref = u.String()
		if s == "#" {
			sn.ref = "#"
		}
		sn.refURL = u

	case "$recursiveRef":
		s, err := value.GetString()
		if err != nil {
			return errWithPath(fmt.Errorf("$recursiveRef must be a string"), sn)
		}

		if s != "#" {
			return errWithPath(fmt.Errorf(
				"$recursiveRef's value must be '#'"), sn)
		}

		sn.recursiveRef = s

	case "$recursiveAnchor":
		b, err := value.GetBool()
		if err != nil {
			return errWithPath(fmt.Errorf(
				"$recursiveAnchor must be a boolean"), sn)
		}

		sn.recursiveAnchor = b

	case "$defs":
		pn := newPathNodeString("$defs")
		node, err := doParseObjectSchema(value, sn, pn)
		if err != nil {
			return err
		}

		sn.kvMap["$defs"] = node
	}
	return nil
}

func loadAssersion(key string, v *jvalue.V, sn *schemaNode) error {
	f, ok := gAssertionConstructorMap[key]
	if !ok {
		return nil
	}

	pn := newPathNodeString(key)
	jn, err := doParseJSON(v, sn, pn)
	if err != nil {
		return err
	}
	ast, err := f(key, jn)
	if err != nil {
		return err
	}
	sn.assertions = append(sn.assertions, ast)
	return nil
}

func loadApplicator(key string, v *jvalue.V, sn *schemaNode) error {
	if key == "items" {
		return loadApplicatorItems(v, sn)
	}

	asf, asfOK := gApplicatorArraySchemaMap[key]
	osf, osfOK := gApplicatorObjectSchemaMap[key]
	sf, sfOK := gApplicatorSchemaMap[key]

	if !asfOK && !osfOK && !sfOK {
		return nil
	}

	pn := newPathNodeString(key)

	var err error
	var node *schemaNode
	var apt applicator

	if asfOK {
		node, err = doParseArraySchema(v, sn, pn)
		if err != nil {
			return err
		}
		apt, err = asf(key, node)
	}

	if osfOK {
		node, err = doParseObjectSchema(v, sn, pn)
		if err != nil {
			return err
		}
		apt, err = osf(key, node)
	}

	if sfOK {
		node, err = doParseSchema(v, sn, pn)
		if err != nil {
			return err
		}
		apt, err = sf(key, node)
	}

	if err != nil {
		return err
	}

	sn.applicators = append(sn.applicators, apt)
	sn.kvMap[key] = node
	return nil
}

func loadApplicatorItems(v *jvalue.V, sn *schemaNode) error {
	pn := newPathNodeString("items")

	var err error
	var node *schemaNode

	if v.JType == jutil.JTArray {
		node, err = doParseArraySchema(v, sn, pn)
	} else {
		node, err = doParseSchema(v, sn, pn)
	}

	if err != nil {
		return err
	}

	apt, err := newApplicatorItems("items", node)
	if err != nil {
		return err
	}

	sn.applicators = append(sn.applicators, apt)
	sn.kvMap["items"] = node
	return nil
}

func finalizeIfThenElse(sn *schemaNode) {
	aptIf := getApplicator(sn, "if", true)
	aptElse := getApplicator(sn, "else", true)
	aptThen := getApplicator(sn, "then", true)

	if aptIf == nil || (aptElse == nil && aptThen == nil) {
		return
	}

	af := aptIf.(*applicatorIf)
	af.elseApt = aptElse
	af.thenApt = aptThen
	sn.applicators = append(sn.applicators, aptIf)
}

func finalizeContains(sn *schemaNode) {
	astMinContains := getAssertion(sn, "minContains", true)
	astMaxContains := getAssertion(sn, "maxContains", true)

	aptContains := getApplicator(sn, "contains", false)
	if aptContains == nil {
		return
	}

	aptc := aptContains.(*applicatorContains)
	if astMinContains != nil {
		ast := astMinContains.(*assertionMinContains)
		aptc.minContains = ast.minContains
	}

	if astMaxContains != nil {
		ast := astMaxContains.(*assertionMaxContains)
		aptc.maxContains = ast.maxContains
	}
}

func finalizeAdditionalItems(sn *schemaNode) {
	aptItems := getApplicator(sn, "items", false)
	aptAdditionalItems := getApplicator(sn, "additionalItems", true)
	if aptItems == nil || aptAdditionalItems == nil {
		return
	}

	aptItemsObj := aptItems.(*applicatorItems)
	if aptItemsObj.node.jtType != jutil.JTArray {
		return
	}

	aptAdditionalItems.(*applicatorAdditionalItems).itemsArrLen =
		len(aptItemsObj.node.nodeArr)
	sn.applicators = append(sn.applicators, aptAdditionalItems)
}

func finalizeAdditionalProperties(sn *schemaNode) {
	aptAdditionalP := getApplicator(sn, "additionalProperties", false)
	if aptAdditionalP == nil {
		return
	}
	aptAP := aptAdditionalP.(*applicatorAdditionalProperties)

	aptProperties := getApplicator(sn, "properties", false)
	if aptProperties != nil {
		aptP := aptProperties.(*applicatorProperties)
		for k := range aptP.node.kvMap {
			aptAP.properties = append(aptAP.properties, k)
		}
	}

	aptPatternProperties := getApplicator(sn, "patternProperties", false)
	if aptPatternProperties != nil {
		aptPP := aptPatternProperties.(*applicatorPatternProperties)
		for _, v := range aptPP.psArr {
			aptAP.patternProperties = append(aptAP.patternProperties, v.reg)
		}
	}
}
