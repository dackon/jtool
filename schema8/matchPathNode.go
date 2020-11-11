package schema8

type matchPathNode struct {
	node          *schemaNode
	path          *pathNode
	assertionName string
}

func newMatchPathNode(s *schemaNode) *matchPathNode {
	if len(s.path) == 0 {
		// Here s must be boolean schema
		if s.snType != sntBoolSchema || s.parent != nil {
			panic("schema must be root boolean schema")
		}
		return &matchPathNode{
			node: s,
		}
	}

	return &matchPathNode{
		node: s,
		path: s.path[len(s.path)-1],
	}
}

func newMatchPathNodeFromString(s *schemaNode, str string) *matchPathNode {
	return &matchPathNode{
		node: s,
		path: newPathNodeString(str),
	}
}
