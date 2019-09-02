package schema

import (
	"fmt"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type definitions struct {
	schemaImpl *schemaImpl

	schemaMap map[string]*Schema
}

func newDefinitions(s *schemaImpl) *definitions {
	return &definitions{
		schemaImpl: s,
		schemaMap:  make(map[string]*Schema),
	}
}

func (def *definitions) loadDefinitions() error {
	jv, ok := def.schemaImpl.jvMap["definitions"]
	if !ok {
		return nil
	}

	if jv.JType != jutil.JTObject {
		return fmt.Errorf("definitions must be object")
	}

	m := jv.Value.(map[string]*jvalue.V)
	for k, v := range m {
		key := objectKey(def.schemaImpl.Schema, "definitions", k)
		s, err := parse(key, v, def.schemaImpl.Schema,
			def.schemaImpl.Schema.root)
		if err != nil {
			return err
		}
		def.schemaMap[k] = s
	}

	def.schemaImpl.svMap["definitions"] = newSchemaVObj(def.schemaMap)
	return nil
}
