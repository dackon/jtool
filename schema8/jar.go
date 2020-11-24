package schema8

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// jar ...
type jar struct {
	sm map[string]*schemaNode
}

// newJar ...
func newJar() *jar {
	return &jar{sm: make(map[string]*schemaNode)}
}

func (j *jar) Add(s *schemaNode) error {
	if len(s.canonicalURIs) == 0 {
		return fmt.Errorf("canonicalURIs is empty")
	}

	for _, u := range s.canonicalURIs {
		j.sm[u] = s
	}

	return nil
}

func (j *jar) AddSchemaJar(schemaJar *jar) error {
	for k, v := range schemaJar.sm {
		j.sm[k] = v
	}
	return nil

}

func (j *jar) Get(uri string) (*schemaNode, error) {
	s, ok := j.sm[uri]
	if ok {
		return s, nil
	}

	u, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("bad URL is %s. Err is %s", uri, err)
	}
	u.Fragment = ""

	_, ok = j.sm[u.String()]
	if ok {
		return nil, fmt.Errorf("failed to find URI %s", uri)
	}

	if gResolverFunc != nil {
		schema, err := gResolverFunc(uri)
		if err == nil {
			s, err := schema.schemaJar.Get(uri)
			if err == nil {
				if err = j.AddSchemaJar(schema.schemaJar); err != nil {
					return nil, err
				}
				return s, nil
			}
		}
	}

	var node *schemaNode

	switch u.Scheme {
	case "http", "https":
		node, err = j.loadProtocolHTTPSchema(u.String())
	case "file":
		node, err = j.loadProtocolFileSchema(u.String())
	default:
		return nil, fmt.Errorf("unknown protocol %s", u.Scheme)
	}

	if err != nil {
		return nil, err
	}

	err = j.AddSchemaJar(node.schema.schemaJar)
	if err != nil {
		return nil, err
	}

	sn, err := node.schema.schemaJar.Get(uri)
	if err != nil {
		return nil, err
	}

	return sn, nil
}

func (j *jar) loadProtocolHTTPSchema(uri string) (*schemaNode, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to download schema. URL is %s. Err is %s", uri, err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body. URL is %s. Err is %s",
			uri, err)
	}

	s, err := doParse(data, uri)
	if err != nil {
		return nil, err
	}
	return s.root, nil
}

func (j *jar) loadProtocolFileSchema(uri string) (*schemaNode, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("bad URL is %s. Err is %s", uri, err)
	}

	if u.Host != "" && u.Host != "127.0.0.1" && u.Host != "localhost" {
		return nil, fmt.Errorf("host can only be empry, '127.0.0.1' or "+
			"'localhost'. URL is %s", uri)
	}

	data, err := ioutil.ReadFile(u.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s. Err is %s", uri, err)
	}

	s, err := doParse(data, uri)
	if err != nil {
		return nil, err
	}
	return s.root, nil
}

func (j *jar) LinkRef() error {
	var err error
	for _, s := range j.sm {
		if s.ref == "" {
			continue
		}

		if s.refURL.IsAbs() {
			s.refSchema, err = j.Get(s.ref)
			if err != nil {
				return err
			}
			continue
		}

		uri := s.baseURIObj.ResolveReference(s.refURL)
		s.refSchema, err = j.Get(uri.String())
		if err != nil {
			return err
		}
	}

	return nil
}

func (j *jar) debug() {
	log.Printf("=========================JAR Contents=========================")
	for k, v := range j.sm {
		log.Printf("Jar::debug:==> Jar is %p. Key is %s. schemaNode is %s.",
			j, k, v)
	}
	log.Printf("=======================JAR Contents End=======================")
}
