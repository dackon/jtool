# Introduction
Schema implemented JSON schema [draft-07](https://json-schema.org/).

**Limitations**: Doesn't support external schema reference.

## Howto
```golang
    tpl := "{\"type\": \"object\"}"
    target := "{\"abc\":1}"
    s := schema.Parse([]byte(tpl))
    if err := schema.Match([]byte(target)); err != nil {
        return errors.New("match failed")
    }
```
