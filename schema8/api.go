package schema8

import "encoding/json"

// RegisterResolver ...
func RegisterResolver(rv ResolverFunc) {
	gResolverFunc = rv
}

// RegisterFormatFunc ...
func RegisterFormatFunc(format string, f FormatFunc) {
	gFormatFuncMap[format] = f
}

// Parse ...
func Parse(raw json.RawMessage, canonicalURL string) (*Schema, error) {
	return doParse(raw, canonicalURL)
}
