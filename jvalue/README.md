# Introduction

JValue (JSON Value) provides support for [schema](../schema), [template](../template) and [operator](../operator) to interact with arbitrary JSON.

## Usage
```golang
    // Parse JSON bytes to JValue.
    v, err := jvalue.ParseJSON([]byte("..."))
    if err != nil {
        return err
    }
    
    // We can get go-value from JValue easily.
    strArr, err := v.GetStringArr()
    if err != nil {
        return err
    }
    
    fmt.Println(strArr)
    
    // Marshal to JSON bytes.
    raw, err := v.Marshal()
    if err != nil {
        return err
    }
    
    fmt.Println(string(raw))
```
