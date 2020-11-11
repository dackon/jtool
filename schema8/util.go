package schema8

import (
	"fmt"
	"net/url"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

func pathToFragment(arr []*pathNode) string {
	fragment := "#"
	for i, s := range arr {
		if i == 0 {
			continue
		}
		tmp := ""
		if s.Type == pntString {
			tmp = url.PathEscape(s.Value.(string))
		} else {
			tmp = fmt.Sprintf("%d", s.Value.(int))
		}
		fragment = fmt.Sprintf("%s/%s", fragment, tmp)
	}

	return fragment
}

func errWithPath(err error, sn *schemaNode) error {
	return fmt.Errorf("err is %s. JSONPointer is %s. node is %s",
		err, pathToFragment(sn.path), sn)
}

func isEqual(s *schemaNode, v *jvalue.V) bool {
	if s.jtType != v.JType {
		return false
	}

	switch s.jtType {
	case jutil.JTNull:
		return true
	case jutil.JTInteger:
		return s.value.(int64) == v.Value.(int64)
	case jutil.JTNumber:
		return jutil.FloatEquals(s.value.(float64), v.Value.(float64))
	case jutil.JTBoolean:
		return s.value.(bool) == v.Value.(bool)
	case jutil.JTString:
		return s.value.(string) == v.Value.(string)
	case jutil.JTObject:
		vvm := v.Value.(map[string]*jvalue.V)
		if len(s.kvMap) != len(vvm) {
			return false
		}

		for key, value := range s.kvMap {
			vv, ok := vvm[key]
			if !ok {
				return false
			}

			if !isEqual(value, vv) {
				return false
			}
		}

		return true
	case jutil.JTArray:
		varr := v.Value.([]*jvalue.V)
		if len(s.nodeArr) != len(varr) {
			return false
		}

		for i, sv := range s.nodeArr {
			if !isEqual(sv, varr[i]) {
				return false
			}
		}
		return true
	}

	return false
}

func getApplicator(sn *schemaNode, name string, remove bool) applicator {
	var i int
	var apt applicator
	for i, apt = range sn.applicators {
		if apt.Name() == name {
			break
		}
		apt = nil
	}
	if apt != nil && remove {
		sn.applicators = append(sn.applicators[0:i],
			sn.applicators[i+1:]...)
	}
	return apt
}

func getAssertion(sn *schemaNode, name string, remove bool) assertion {
	var i int
	var ast assertion
	for i, ast = range sn.assertions {
		if ast.Name() == name {
			break
		}
		ast = nil
	}
	if ast != nil {
		sn.assertions = append(sn.assertions[0:i], sn.assertions[i+1:]...)
	}
	return ast
}

func copyPath(sn *schemaNode, pn *pathNode) []*pathNode {
	pathArr := make([]*pathNode, len(sn.path)+1)
	for i, p := range sn.path {
		pathArr[i] = p
	}
	pathArr[len(sn.path)] = pn
	return pathArr
}
