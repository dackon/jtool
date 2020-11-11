package schema8

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"time"

	"github.com/dackon/jtool/jutil"
	"github.com/dackon/jtool/jvalue"
)

var (
	gEmailReg    = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

func formatMySQLDateTime(v *jvalue.V) error {
	if v.JType != jutil.JTString {
		return nil
	}

	s, _ := v.GetString()
	f := "2006-01-02 15:04:05"
	if _, err := time.Parse(f, s); err == nil {
		return nil
	}

	return fmt.Errorf("failed to parse '%s' to DateTime", s)
}

func formatMongoDBDateTime(v *jvalue.V) error {
	if v.JType != jutil.JTString {
		return nil
	}

	s, _ := v.GetString()
	if _, err := time.Parse(time.RFC3339Nano, s); err == nil {
		return nil
	}

	return fmt.Errorf("failed to parse '%s' to DateTime", s)
}

func formatEmail(v *jvalue.V) error {
	if v.JType != jutil.JTString {
		return nil
	}

	s, _ := v.GetString()
	if gEmailReg.MatchString(s) {
		return nil
	}

	return fmt.Errorf("bad email address '%s'", s)
}

func formatIPV4(v *jvalue.V) error {
	if v.JType != jutil.JTString {
		return nil
	}

	s, _ := v.GetString()
	ipv4 := net.ParseIP(s)
	if ipv4 == nil {
		return fmt.Errorf("bad ip address '%s'", s)
	}

	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '.':
			return nil
		case ':':
			return fmt.Errorf("ip address '%s' is ipv6 not ipv4", s)
		}
	}

	return nil
}

func formatIPV6(v *jvalue.V) error {
	if v.JType != jutil.JTString {
		return nil
	}

	s, _ := v.GetString()
	ipv6 := net.ParseIP(s)
	if ipv6 == nil {
		return fmt.Errorf("bad ip address '%s'", s)
	}

	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '.':
			return fmt.Errorf("ip address '%s 'is ipv4 not ipv6", s)
		case ':':
			return nil
		}
	}

	return nil
}

func formatURI(v *jvalue.V) error {
	if v.JType != jutil.JTString {
		return nil
	}

	s, _ := v.GetString()
	u, err := url.Parse(s)
	if err != nil {
		return fmt.Errorf("failed to parse URL '%s', err is '%s'", s, err)
	}

	if !u.IsAbs() {
		return fmt.Errorf("URL '%s' is not absolute URL", s)
	}

	return nil
}

func formatURIRef(v *jvalue.V) error {
	if v.JType != jutil.JTString {
		return nil
	}

	s, _ := v.GetString()
	_, err := url.Parse(s)
	if err != nil {
		return fmt.Errorf("failed to parse URL '%s', err is '%s'", s, err)
	}

	return nil
}

func formatRegex(v *jvalue.V) error {
	if v.JType != jutil.JTString {
		return nil
	}

	s, _ := v.GetString()
	_, err := regexp.Compile(s)
	if err != nil {
		return fmt.Errorf("failed to parse '%s' to regexp, err is '%s'", s, err)
	}

	return nil
}

func formatTimestamp(v *jvalue.V) error {
	if v.JType != jutil.JTInteger {
		return nil
	}

	i, _ := v.GetInteger()
	now := time.Now().Unix()

	if i < now-600 || i > now+600 {
		return fmt.Errorf("Current timestamp is %d. Value is %d", now, i)
	}

	return nil
}

func formatTimestampMS(v *jvalue.V) error {
	if v.JType != jutil.JTInteger {
		return nil
	}

	i, _ := v.GetInteger()
	now := time.Now().UnixNano()
	now /= 1000000

	if i < now-600000 || i > now+600000 {
		return fmt.Errorf("Current timestamp is %d. Value is %d", now, i)
	}

	return nil
}

func formatUUID(v *jvalue.V) error {
	if v.JType != jutil.JTInteger {
		return nil
	}

	s, _ := v.GetString()
	_, err := parseUUID(s)
	if err != nil {
		return fmt.Errorf("'%s' is not a valid UUID", s)
	}

	return nil
}
