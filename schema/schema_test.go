package schema

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func doTest(jname, sname string, t *testing.T) {
	data, err := ioutil.ReadFile(jname)
	if err != nil {
		t.Fatal(err)
	}
	sdata, err := ioutil.ReadFile(sname)
	if err != nil {
		t.Fatal(err)
	}

	schema, err := Parse(sdata)
	if err != nil {
		fmt.Println(jname, sname)
		t.Fatal(err)
	}

	err = schema.Match(data)
	if err != nil {
		t.Fatalf("schema file is %s. failed to match. Err is %s.", sname, err)
	}
}

func TestArray(t *testing.T) {
	doTest("./testdata/array.json", "./testdata/arraySchema.json", t)
}

func TestBoolLogic(t *testing.T) {
	doTest("./testdata/boolLogic.json", "./testdata/boolLogicSchema.json", t)
}

func TestConditions(t *testing.T) {
	doTest("./testdata/conditions.json", "./testdata/conditionsSchema.json", t)
}

func TestConst(t *testing.T) {
	doTest("./testdata/const.json", "./testdata/constSchema.json", t)
}

func TestEnum(t *testing.T) {
	doTest("./testdata/enum.json", "./testdata/enumSchema.json", t)
}

func TestNumber(t *testing.T) {
	doTest("./testdata/number.json", "./testdata/numberSchema.json", t)
}

func TestObject(t *testing.T) {
	doTest("./testdata/object.json", "./testdata/objectSchema.json", t)
}

func TestRef(t *testing.T) {
	doTest("./testdata/ref.json", "./testdata/refSchema.json", t)
}

func TestString(t *testing.T) {
	doTest("./testdata/string.json", "./testdata/stringSchema.json", t)
}

func TestType(t *testing.T) {
	doTest("./testdata/type.json", "./testdata/typeSchema.json", t)
}

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

func TestAll(t *testing.T) {
	path := "./test-suite"
	dir, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}

	fileInfoArr, err := dir.Readdir(1000)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < len(fileInfoArr); i++ {
		if fileInfoArr[i].IsDir() {
			continue
		}

		fpath := path + string(os.PathSeparator) + fileInfoArr[i].Name()
		if fpath[len(fpath)-1] == '~' {
			os.Remove(fpath)
			continue
		}
		doTestSingle(path+string(os.PathSeparator)+fileInfoArr[i].Name(), t)
	}
}

func doTestSingle(path string, t *testing.T) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	//	t.Logf("Path is %s", path)

	var testArr []testItems
	err = json.Unmarshal(data, &testArr)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < len(testArr); i++ {
		schema, err := Parse(testArr[i].Schema)
		if err != nil {
			t.Fatalf("Compile failed. Err is %s. Path is %s. "+
				"Index is %d. Schema is %s.", err, path, i, testArr[i].Schema)
		}

		for j := 0; j < len(testArr[i].Tests); j++ {
			ti := testArr[i].Tests[j]
			//			t.Logf("Path is %s. Schema is %s. Test is %s",
			//				path, testArr[i].Schema, ti.Data)
			err = schema.Match(ti.Data)
			if (err != nil && ti.Valid) || (err == nil && !ti.Valid) {
				t.Fatalf("Validator failed. Path is %s. Schema is %s. "+
					"Test is %s. i is %d. j is %d. Err is %s.",
					path, testArr[i].Schema, ti.Data, i, j, err)
			}
		}
	}

	t.Logf("%s passed test", path)
}
