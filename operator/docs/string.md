# Operations for String
1. [$concat](#concat)
2. [$sha256](#sha256)
3. [$toLower](#toLower)
4. [$toUpper](#toUpper)


# $concat
$concat concatenates items in array to string. If item in array is not string, it will be converted to string.

For example, giving JSON:
```json
{
    "a": {
        "$concat": ["abc", "eft", ""]
    },
    "b": {
        "$concat": ["abc", 1, true, null]
    }
}
```
will produce:
```json
{
    "a": "abceft",
    "b": "abc1truenull"
}
```

# $sha256
$sha256 computes the sha256 of the giving string.

For example, giving JSON:
```json
{
    "a": {
        "$sha256": "abc"
    }
}
```
will produce:
```json
{
    "a": "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"
}
```

# $toLower
$toLower gives the lower case of the string.

For example, giving JSON:
```json
{
    "a": {
        "toLower": "aAB"
    }
}
```
will produce:
```json
{
    "a": "aab"
}
```

# $toUpper
$toUpper gives the upper case of the string.

For example, giving JSON:
```json
{
    "a": {
        "$toUpper": "aAB"
    }
}
```
will produce:
```json
{
    "a": "AAB"
}
```
