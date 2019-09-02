# Table of Contents
1. [Introduction](#Introduction)
2. [Example](#Example)
3. [Usage](#Usage)
4. [Array Operators](docs/array.md)
5. [Date Operators](docs/date.md)
6. [Number Operators](docs/number.md)
7. [String Operators](docs/string.md)
8. [Playground](#Playground)

# Introduction
This package implemented lots of operators to manipulate JSON values. Currently, only a few operators were implemented.

## Example
Giving JSON:
```json
{
    "total": {
        "$sum": [1, 2]
    }
}
```

Here, '$sum' will compute the sum of the array '[1, 2]', so the result will be:
```json
{
    "total": 3
}
```

## Usage
```golang
    /* 
        // Register your own operators if you have any.
        err := operator.RegisterOP("$sha512", myOP)
        if err != nil {
            return err
        }
    */

    // First, parse JSON bytes to JValue.
    jv, err := jvalue.ParseJSON([]byte(...))
    if err != nil {
        return err
    }

    // Do the operation.
    newV, err := operator.Do(jv)
    if err != nil {
        return err
    }
    
    raw, err := newV.Marshal()
    fmt.Println(raw)
```

## Playground
[Goto GO Playground](https://play.golang.org/p/GWHASNc_BhN)
