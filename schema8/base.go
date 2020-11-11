package schema8

type base struct {
	name string
	node *schemaNode
}

func (b *base) Name() string {
	return b.name
}

func (b *base) GetSchemaNode() *schemaNode {
	return b.node
}
