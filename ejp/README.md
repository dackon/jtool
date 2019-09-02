# Introduction
This project implements 'Extended JSON Pointer' (EJP).

## Differences from RFC6901
- Each EJP has a name of the JSON it references to, the charset for name is [\_a-zA-Z0-9].
- In RFC6901, '\~' must be encoded to '\~0', '/' must be encoded to '\~1'. Here, besides '\~' and '/', the '$' must be encoded to '\~2', which is used by [JSON template](../template) to reference each item in array.
- '-' specifies the last item in array.

## Examples
Giving the following two JSONs:
JSON name: foo
```json
{
  "a": 1,
  "b": {
    "b1": "apple"
  },
  "c": [4, 5, 6]
}
```

JSON Name: bar
```json
{
  "x": [{"x1":10}, {"x1":20}, {"x1": 30}],
  "y": "world"
}
```

Referenced Values

| EJP        | Referenced Value                               |
|------------|------------------------------------------------|
| foo/a      | 1                                              |
| foo/b      | {"b1": "apple"}                                |
| foo/b/b1   | "apple"                                        |
| foo/c/1    | 5                                              |
| foo/c/-    | 6                                              |
| bar/x/0    | {"x1":10}                                      |
| bar/x/1/x1 | 20                                             |
| bar/y      | "world"                                        |
| foo        | {"a": 1, "b": {"b1": "apple"}, "c": [4, 5, 6]} |
