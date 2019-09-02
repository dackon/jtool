package operator

import (
	"fmt"
	"regexp"
	"time"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

type dateAddMS struct {
	reg *regexp.Regexp
}

func newDateAddMS() *dateAddMS {
	return &dateAddMS{reg: regexp.MustCompile(
		`^[0-9]{4}\-[0-9]{2}\-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2}$`)}
}

// do will add milliseconds to the date, e.g.:
// {"s": {"$dateAddMS":["2011-01-01 11:11:00", 1000]}} -> {"s": "2011-01-01 11:11:01"}
// {"s": {"$dateAddMS":["2019-06-18T07:21:48Z", 1000]}} -> {"s": "2019-06-18T07:21:49Z"}
func (op *dateAddMS) do(jv *jvalue.V) (*jvalue.V, error) {
	if jv.JType != jutil.JTArray {
		return nil, fmt.Errorf("operator '$dateAddMS': value is not array")
	}

	arrv := jv.Value.([]*jvalue.V)
	if len(arrv) != 2 || arrv[0].JType != jutil.JTString ||
		arrv[1].JType != jutil.JTInteger {
		return nil, fmt.Errorf("operator '$dateAddMS': bad param")
	}

	e1, _ := arrv[0].GetString()
	e2, _ := arrv[1].GetInteger()

	var t time.Time
	var err error
	var tstr string
	var isRFC3339 bool

	if op.reg.MatchString(e1) {
		isRFC3339 = false
		t, err = time.Parse("2006-01-02 15:04:05", e1)
	} else {
		isRFC3339 = true
		t, err = time.Parse(time.RFC3339, e1)
	}

	if err != nil {
		return nil, fmt.Errorf("operator '$dateAddMS': parse %s failed", e1)
	}
	t = t.Add(time.Millisecond * time.Duration(e2))

	if !isRFC3339 {
		tstr = t.Format("2006-01-02 15:04:05")
	} else {
		tstr = t.Format(time.RFC3339)
	}

	newJV := &jvalue.V{JType: jutil.JTString, Value: tstr}
	return newJV, nil
}
