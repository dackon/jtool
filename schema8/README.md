# Introduction
This project implemented JSON schema [draft 2019-09 (draft-08)](https://json-schema.org/).

# Limitations
1. Doesn't support meta-shema, '$schema' and '$vocabulary'
2. Doesn't support collecting annotations, applicator 'unevaluatedProperties' and 'unevaluatedItems'
3. Doesn't support vocabulary for the contents of string-encoding data
4. Doesn't support vocabulary for basic meta-data annotations
5. Doesn't support match length of unicode code point string
6. You must use the normalized URL, this implementation will **NOT** normalize it for you
7. 'format' is assertion **ONLY**, you can register your own format to add new (or replace default) handler
8. Supported output format: Flag (with the error location: keywordLocation, instanceLocation and absoluteKeywordLocation)
9. Only a few formats are suppored. see 'init.go' for details.

# Other important things
1. The default URL of schema is 'https://default.uri' if you do not have $id field in root schema.
2. This implementation treats 0 as integer, 0.0 as float64, and 0 != 0.0
3. Support http, https and file protocols to reference external schema. The host of file protocol URL must be 'localhost', '127.0.0.1' or empty string ''. You can also register your own schema resolver, see the example in file 'schema_test.go'.

# Howto
```golang
    tpl := "{\"type\": \"object\"}"
    target := "{\"abc\":1}"
    s := schema8.Parse([]byte(tpl), "")
    if ret := schema.Match([]byte(target)); ret.Err != nil {
        return fmt.Errorf("schema.Match failed. Result is %s", ret)
    }
```
