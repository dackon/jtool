package schema8

import (
	"context"
	"encoding/json"

	"github.com/dackon/jtool/jvalue"
)

type Schema struct {
	root      *schemaNode
	schemaJar *jar
}

func (s *Schema) Match(raw json.RawMessage) *MatchInfo {
	v, err := jvalue.ParseJSON(raw)
	if err != nil {
		return &MatchInfo{Err: err}
	}

	mc := &matchContext{mi: &MatchInfo{}}
	ctx := context.WithValue(context.Background(), matchContextKey, mc)
	return s.root.match(ctx, v)
}
