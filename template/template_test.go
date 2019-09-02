package template

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"github.com/dackon/jtool/jvalue"
	"github.com/dackon/jtool/schema"
)

type dataSourceCfg struct {
	Name string          `json:"name"`
	Data json.RawMessage `json:"data"`
}

type oneCaseCfg struct {
	Case     string          `json:"case"`
	Enable   bool            `json:"enable"`
	JPool    []dataSourceCfg `json:"jpool"`
	Template json.RawMessage `json:"template"`
	Schema   json.RawMessage `json:"schema"`
}

type oneCase struct {
	Case   string
	JPool  *JSONPool
	Tpl    *Template
	Schema *schema.Schema
}

func loadTestCase() ([]*oneCase, error) {
	path := "./testcase/cases.json"

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	caseCfgArr := []*oneCaseCfg{}
	if err = json.Unmarshal(data, &caseCfgArr); err != nil {
		return nil, err
	}

	cases := make([]*oneCase, 0, len(caseCfgArr))
	for i, v := range caseCfgArr {
		if !v.Enable {
			continue
		}

		c := &oneCase{
			Case:  v.Case,
			JPool: NewJSONPool(),
		}

		for _, d := range v.JPool {
			j, err := jvalue.ParseJSON(d.Data)
			if err != nil {
				return nil, fmt.Errorf("loadTestCase: ParseJSON failed. "+
					"idx is %d, err is %s", i, err)
			}
			c.JPool.Add(d.Name, j)
		}

		c.Tpl, err = Parse(v.Template)
		if err != nil {
			return nil, fmt.Errorf("loadTestCase: Template Compile failed. "+
				"idx is %d, err is %s", i, err)
		}

		c.Schema, err = schema.Parse(v.Schema)
		if err != nil {
			return nil, fmt.Errorf("loadTestCase: Schema Compile failed. "+
				"idx is %d, err is %s", i, err)
		}
		cases = append(cases, c)
	}

	return cases, nil
}

func doTest(caseArr []*oneCase) error {
	for _, c := range caseArr {
		jv, err := c.Tpl.Execute(c.JPool)
		if err != nil {
			return fmt.Errorf("Case %s Execute failed. Err is %s", c.Case, err)
		}

		raw, err := jv.Marshal()
		if err != nil {
			return fmt.Errorf("Case %s Marshal failed. Err is %s", c.Case, err)
		}

		if err = c.Schema.Match(raw); err != nil {
			return fmt.Errorf("Case %s Match failed. Template generated JSON "+
				"is %s. Err is %s", c.Case, string(raw), err)
		}
	}

	return nil
}

func TestTemplate(t *testing.T) {
	cases, err := loadTestCase()
	if err != nil {
		log.Fatalf("loadTestCase error %s", err)
	}

	err = doTest(cases)
	if err != nil {
		log.Fatalf("doTest failed. Err is %s", err)
	}
}
