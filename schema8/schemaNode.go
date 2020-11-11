package schema8

import (
	"context"
	"fmt"
	"net/url"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type schemaNode struct {
	snType schemaNodeType

	jtType jutil.JType
	value  interface{}

	kvMap   map[string]*schemaNode
	nodeArr []*schemaNode

	// the '$schema'
	metaSchema string

	schema *Schema

	// baseURI
	baseURI       string
	baseURIObj    *url.URL
	canonicalURIs []string

	// $id
	id     string
	anchor string
	idURL  *url.URL

	// $ref
	ref       string
	refURL    *url.URL
	refSchema *schemaNode

	recursiveRef    string
	recursiveAnchor bool

	format string

	parent *schemaNode
	path   []*pathNode

	assertions  []assertion
	applicators []applicator
}

// finalize setBaseURI and add schemaNode to the jar
func (s *schemaNode) finalize() (err error) {
	if s.snType != sntNotSchema {
		if s.id != "" {
			if s.idURL.IsAbs() {
				s.baseURI = s.id
				s.baseURIObj = s.idURL
			} else {
				node := s.parent
				for node != nil {
					if node.baseURI != "" {
						s.baseURIObj = node.baseURIObj.ResolveReference(s.idURL)
						s.baseURI = s.baseURIObj.String()
						break
					} else {
						node = node.parent
					}
				}

				if s.baseURI == "" {
					panic("baseURI must not be empty")
				}
			}
		} else {
			node := s.parent
			for node != nil {
				if node.baseURI != "" {
					s.baseURI = node.baseURI
					s.baseURIObj = node.baseURIObj
					break
				}
				node = node.parent
			}
		}

		s.setCanonicalURI()
		err = s.schema.schemaJar.Add(s)
		if err != nil {
			return err
		}
	}

	// If here, s must have baseURI
	for _, v := range s.kvMap {
		if err = v.finalize(); err != nil {
			return err
		}
	}

	for _, v := range s.nodeArr {
		if err = v.finalize(); err != nil {
			return err
		}
	}

	return nil
}

func (s *schemaNode) setCanonicalURI() error {
	if s.anchor != "" {
		cu := fmt.Sprintf("%s#%s", s.baseURI, s.anchor)
		s.canonicalURIs = append(s.canonicalURIs, cu)
	}

	node := s
	var arr []*pathNode
	for node != nil {
		if node.isRoot() {
			arr = append(arr, newPathNodeString(""))
			for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
				arr[i], arr[j] = arr[j], arr[i]
			}

			jp := pathToFragment(arr)
			cu := fmt.Sprintf("%s%s", s.baseURI, jp)
			s.canonicalURIs = append(s.canonicalURIs, cu)
			break
		}

		arr = append(arr, node.path[len(node.path)-1])
		node = node.parent
	}

	for i := 0; i < len(s.canonicalURIs); i++ {
		obj, err := url.Parse(s.canonicalURIs[i])
		if err != nil {
			return errWithPath(err, s)
		}
		s.canonicalURIs[i] = obj.String()
	}

	return nil
}

func (s *schemaNode) match(ctx context.Context, jv *jvalue.V) *MatchInfo {
	mc := ctx.Value(matchContextKey).(*matchContext)
	copyPath := make([]*matchPathNode, len(mc.path))
	copy(copyPath, mc.path)
	mc.path = append(mc.path, newMatchPathNode(s))
	defer func() {
		mc.path = copyPath
	}()

	return s.doMatch(ctx, jv)
}

func (s *schemaNode) doMatch(ctx context.Context, jv *jvalue.V) *MatchInfo {
	mc := ctx.Value(matchContextKey).(*matchContext)
	switch s.snType {
	case sntBoolSchema:
		if !s.value.(bool) {
			err := fmt.Errorf("boolean schema match failed")
			setFailedMatchInfo(s, jv, err, ctx)
		}
		return mc.mi
	case sntObjectSchema:
		if s.refSchema != nil {
			newmp := newMatchPathNodeFromString(s.refSchema, "$ref")
			mc.path = append(mc.path, newmp)
			if s.refSchema.doMatch(ctx, jv).Err != nil {
				return mc.mi
			}
		}

		if s.recursiveRef != "" {
			rs := s.getRecursiveRefSchema(mc.path)
			newmp := newMatchPathNodeFromString(rs, "$recursiveRef")
			mc.path = append(mc.path, newmp)
			if rs.refSchema.doMatch(ctx, jv).Err != nil {
				return mc.mi
			}
		}

		var err error
		for _, ast := range s.assertions {
			if err = ast.Valid(ctx, jv); err != nil {
				mc.path[len(mc.path)-1].assertionName = ast.Name()
				setFailedMatchInfo(ast.GetSchemaNode(), jv, err, ctx)
				return mc.mi
			}
		}

		for _, apt := range s.applicators {
			if err = apt.Valid(ctx, jv); err != nil {
				return mc.mi
			}
		}
		return mc.mi
	}

	panic("must not be here")
}

func (s *schemaNode) getRecursiveRefSchema(arr []*matchPathNode) *schemaNode {
	if !s.recursiveAnchor {
		cu := s.baseURI + "#"
		sn, err := s.schema.schemaJar.Get(cu)
		if err != nil {
			panic(err)
		}

		return sn
	}

	for _, n := range arr {
		if !n.node.recursiveAnchor {
			continue
		}

		cu := s.baseURI + "#"
		sn, err := s.schema.schemaJar.Get(cu)
		if err != nil {
			panic(err)
		}

		return sn
	}

	panic("must not be here")
}

func (s *schemaNode) isRoot() bool {
	if s.id != "" {
		return true
	}

	return false
}

func (s *schemaNode) String() string {
	keys := make([]string, 0, len(s.kvMap))
	for k := range s.kvMap {
		keys = append(keys, k)
	}

	asts := make([]string, 0, len(s.assertions))
	for _, v := range s.assertions {
		asts = append(asts, v.Name())
	}

	apts := make([]string, 0, len(s.applicators))
	for _, v := range s.applicators {
		apts = append(apts, v.Name())
	}

	return fmt.Sprintf("<<Type: '%s'. JType: '%s'. BaseURI: '%s'. "+
		"CanonicalURIs: %v. Anchor: '%s'. $id: '%s'. Keys: %v. "+
		"ArrayLen: %d. Assertions: %v. Applicators: %v. $ref: '%s'. "+
		"$recursiveRef:'%s'. "+
		"$recursiveAnchor %v. PathArr: %v>>",
		s.snType, s.jtType, s.baseURI, s.canonicalURIs, s.anchor, s.id, keys,
		len(s.nodeArr), asts, apts, s.ref, s.recursiveRef, s.recursiveAnchor,
		s.path)
}
