package operator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/dackon/jtool/jvalue"
)

type caseDataCfg struct {
	Data     json.RawMessage `json:"data"`
	Expected json.RawMessage `json:"expected"` // Optional.
	Error    bool            `json:"error"`    // Optional.
}

type testCaseCfg struct {
	Name   string        `json:"name"`
	Enable bool          `json:"enable"`
	Cases  []caseDataCfg `json:"cases"`
}

type caseData struct {
	Data     *jvalue.V
	Expected *jvalue.V
	Error    bool
}

type testCase struct {
	Name    string
	CaseArr []*caseData
}

func loadOneCategory(name string, t *testing.T) []*testCase {
	fileData, err := ioutil.ReadFile(name)
	if err != nil {
		t.Errorf("Failed to read test case file '%s'.", name)
		return nil
	}

	caseCfgArr := []*testCaseCfg{}
	if err = json.Unmarshal(fileData, &caseCfgArr); err != nil {
		t.Fatalf("Unmarshal file data filed. Err is %s", err)
	}

	caseArr := make([]*testCase, 0, len(caseCfgArr))

	for _, cfg := range caseCfgArr {
		if !cfg.Enable {
			continue
		}

		tc := &testCase{
			Name: cfg.Name,
		}

		for _, c := range cfg.Cases {
			datav, err := jvalue.ParseJSON(c.Data)
			if err != nil {
				t.Fatalf("Data ParseJSON failed. Raw is %s", string(c.Data))
			}

			cd := &caseData{
				Data:  datav,
				Error: c.Error,
			}

			if len(c.Expected) > 0 {
				expectedV, err := jvalue.ParseJSON(c.Expected)
				if err != nil {
					t.Fatalf("Expected ParseJSON failed. Raw is %s",
						string(c.Expected))
				}
				cd.Expected = expectedV
			}
			tc.CaseArr = append(tc.CaseArr, cd)
		}

		caseArr = append(caseArr, tc)
	}

	return caseArr
}

func loadCases(t *testing.T) []*testCase {
	path := "./testCase"
	files, err := ioutil.ReadDir(path)
	if err != nil {
		t.Fatalf("ReadDir failed. Err is %s", err)
	}

	fileNames := []string{}
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		if filepath.Ext(f.Name()) == ".json" {
			fileNames = append(fileNames, f.Name())
		}
	}

	var caseArr []*testCase

	for _, f := range fileNames {
		arr := loadOneCategory(
			fmt.Sprintf("%s%c%s", path, os.PathSeparator, f), t)
		if len(arr) == 0 {
			continue
		}
		caseArr = append(caseArr, arr...)
	}

	return caseArr
}

func TestDo(t *testing.T) {
	caseArr := loadCases(t)

	for _, c := range caseArr {
		for _, d := range c.CaseArr {
			jv, err := Do(d.Data)
			if err != nil && !d.Error {
				t.Fatalf("Case '%s' failed. Expected no error but err is %s.",
					c.Name, err)
			}

			if err == nil && d.Error {
				t.Fatalf("Case '%s' failed. Expected error but err is nil.",
					c.Name)
			}

			if d.Expected != nil {
				if err = jv.IsEqual(d.Expected); err != nil {
					r, _ := jv.Marshal()
					e, _ := d.Expected.Marshal()

					t.Fatalf("Case '%s' failed. Result is '%s'. "+
						"Expected is '%s'.", c.Name, string(r), string(e))
				}
			}
		}

		t.Logf("Success. Case is '%s'.", c.Name)
	}
}
