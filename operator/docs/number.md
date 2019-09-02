# Operations for Number
1. [$abs](#abs)
2. [$ceil](#ceil)
3. [$divide](#divide)
4. [$sum](#sum)

# $abs
$abs computes absolute value of the giving number.

For example, giving JSON:
```json
{
    "a": {
        "$abs": -15.123
    },
    "b": {
        "$abs": 13
    }
}
```
will produce:
```json
{
    "a": 15.123,
    "b": 13
}
```

# $ceil
$ceil computes the least integer value greater than or equal to the giving number.

For example, giving JSON:
```json
{
    "a": {
        "$ceil": 7.2
    },
    "b": {
        "$ceil": -7.2
    }
}
```
will produce:
```json
{
    "a": 8,
    "b": -6
}
```

# $divide
$divide divides the giving number. The value of $divide must be an array, the 1st item in the array is diviend, the 2nd item is divisor.

For example, giving JSON:
```json
{
    "a": {
        "$divide": [1, 3]
    }
}
```
will produce:
```json
{
    "a": 0.3333333333333
}
```

# $sum
$sum computes the sum of the giving array.

For example, giving JSON:
```json
{
    "float": {
        "$sum": [1, 15.23521, -90]
    },
    "int": {
        "$sum": [1000, 12038, 3]
    }
}
```
will produce:
```json
{
    "float": -73.76479,
    "int": 13041
}
```
