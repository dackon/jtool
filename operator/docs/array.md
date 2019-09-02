# Operations for Array
1. [$loop](#loop)

# $loop
$loop applies specified operation to each item in the array. The value of $loop must be an array with two items in it, the first item must be an array (parameter array), the second must be a string specifing the operation. Each item in the parameter array must be suitable for the parameter of the operation.

For example, giving JSON:
```json
{
    "scores": {
        "$loop": [[1, 100, -2], "$abs"]
    }
}
```
will produce:
```json
{
    "scores": [1, 100, 2]
}
```
