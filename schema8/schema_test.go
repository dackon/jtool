package schema8

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"regexp"
	"testing"
)

type testData struct {
	Description string          `json:"description"`
	Data        json.RawMessage `json:"data"`
	Valid       bool            `json:"valid"`
}

type testItems struct {
	Description string          `json:"description"`
	Schema      json.RawMessage `json:"schema"`
	Tests       []testData      `json:"tests"`
}

func resolver(uri string) (*Schema, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	if u.Scheme != "http" && u.Host != "localhost:1234" {
		return nil, fmt.Errorf("unknown uri")
	}

	path := "./tests/remotes" + u.Path

	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema at %s", uri)
	}

	s, err := Parse(f, uri)
	return s, err
}

func TestAll(t *testing.T) {
	RegisterResolver(resolver)

	if testSingle(t) {
		return
	}

	doTestDir("./tests/deep", t)
	doTestDir("./tests/suite", t)
}

func testSingle(t *testing.T) bool {
	fpath := "./tests/single.json"

	if _, err := os.Stat(fpath); err != nil {
		return false
	}

	doTestSingle(fpath, 1, 1, t)
	return true
}

func doTestDir(path string, t *testing.T) {
	dir, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}

	fileInfoArr, err := dir.Readdir(1000)
	if err != nil {
		t.Fatal(err)
	}

	reg := regexp.MustCompile("\\.json$")

	var filteredArr []os.FileInfo
	for i := 0; i < len(fileInfoArr); i++ {
		if fileInfoArr[i].IsDir() {
			continue
		}

		fpath := path + string(os.PathSeparator) + fileInfoArr[i].Name()
		if !reg.MatchString(fpath) {
			continue
		}
		filteredArr = append(filteredArr, fileInfoArr[i])
	}

	for i, f := range filteredArr {
		doTestSingle(path+string(os.PathSeparator)+f.Name(),
			i+1, len(filteredArr), t)
	}
}

func doTestSingle(path string, idx, total int, t *testing.T) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile failed. Err is %s. Path is %s", err, path)
	}

	var testArr []testItems
	err = json.Unmarshal(data, &testArr)
	if err != nil {
		t.Fatalf("Unmarshal failed. Err is %s. Path is %s", err, path)
	}

	for i := 0; i < len(testArr); i++ {
		schema, err := Parse(testArr[i].Schema, "")
		if err != nil {
			t.Fatalf("Parse schema failed. Err is %s. Path is %s. "+
				"Index is %d. Schema is %s.", err, path, i, testArr[i].Schema)
		}

		for j := 0; j < len(testArr[i].Tests); j++ {
			ti := testArr[i].Tests[j]
			//			t.Logf("Path is %s. Schema is %s. Test is %s",
			//				path, testArr[i].Schema, ti.Data)
			//			t.Logf("Next test file is %s. Schema is %s. Schema Idx is %d. "+
			//				"Test case idx is %d. Test JSON value is %s.",
			//	path, testArr[i].Schema, i, j, ti.Data)
			info := schema.Match(ti.Data)
			if (info.Err != nil && ti.Valid) || (info.Err == nil && !ti.Valid) {
				t.Fatalf("Validator failed. Path is %s. Schema is %s. "+
					"Test JSON value is %s. Schema idx is %d. Test case idx "+
					"is %d. Info is %s.",
					path, testArr[i].Schema, ti.Data, i, j, info)
			}
		}
	}

	t.Logf("[%d/%d] %s passed test !!!", idx, total, path)
}
