package schema8

import (
	"encoding/json"
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

func NewJar() *jar {
	return &jar{sm: make(map[string]*schemaNode)}
}

func (j *jar) Parse(raw json.RawMessage, canonicalURL string) (
	*Schema, error) {
	s, err := doParse(raw, canonicalURL)
	if err != nil {
		return nil, err
	}

	if err = j.Add(s.root); err != nil {
		return nil, err
	}

	return s, nil
}

func (j *jar) Add(s *schemaNode) error {
	if len(s.canonicalURIs) == 0 {
		return fmt.Errorf("canonicalURIs is empty")
	}

	for _, u := range s.canonicalURIs {
		_, ok := j.sm[u]
		if ok {
			return errWithPath(fmt.Errorf("duplicated canonicalURI %s", u), s)
		}
		j.sm[u] = s
	}

	return nil
}

func (j *jar) Get(uri string) (*schemaNode, error) {
	s, ok := j.sm[uri]
	if ok {
		return s, nil
	}

	if gResolverFunc != nil {
		schema, err := gResolverFunc(uri)
		if err == nil {
			s, err := schema.schemaJar.Get(uri)
			if err == nil {
				for k, v := range schema.schemaJar.sm {
					_, ok := j.sm[k]
					if ok {
						return nil, errWithPath(fmt.Errorf(
							"duplicated canonicalURI %s", k), v)
					}
					j.sm[k] = v
				}
				return s, nil
			}
		}
	}

	u, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("bad URL is %s. Err is %s", uri, err)
	}

	var node *schemaNode

	fragment := u.Fragment
	u.Fragment = ""

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

	err = j.Add(node)
	if err != nil {
		return nil, err
	}

	sn, err := node.schema.schemaJar.Get(uri)
	if err != nil {
		return nil, err
	}

	if fragment != "" {
		j.Add(sn)
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
		log.Printf("Jar::debug: Key is %s. schemaNode is %s.", k, v)
	}
	log.Printf("=======================JAR Contents End=======================")
}
