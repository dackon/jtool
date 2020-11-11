package schema8

import (
	"context"

	"github.com/dackon/jtool/jvalue"
)

// FormatFunc ...
type FormatFunc func(v *jvalue.V) error

// ResolverFunc
type ResolverFunc func(uri string) (*Schema, error)

type schemaNodeType string
type pathNodeType string
type contextKeyType string

type assertionConstructorFunc func(name string, n *schemaNode) (
	assertion, error)
type applicatorConstructorFunc func(name string, n *schemaNode) (
	applicator, error)

type matchContext struct {
	mi   *MatchInfo
	path []*matchPathNode
}

type assertion interface {
	Name() string
	Valid(ctx context.Context, v *jvalue.V) error
	GetSchemaNode() *schemaNode
}

type applicator interface {
	Name() string
	Valid(ctx context.Context, v *jvalue.V) error
	Annotation() interface{}
}
