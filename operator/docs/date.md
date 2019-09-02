# Operations for Date
1. [$dateAddMS](#dateAddMS)

# $dateAddMS
$dateAddMS adds milliseconds to date string. It supports two formats date string: 'YYYY-MM-DD HH:mm:ss' and RFC3339.

For example, giving JSON:
```json
{
    "date1": {
        "$dateAddMS": ["2011-11-11 00:00:00", 1000]
    },
    "date2": {
        "$dateAddMS": ["2019-06-10T06:48:50Z", 1000]
    }
}
```
will produce:
```json
{
    "date1": "2011-11-11 00:00:01",
    "date2": "2019-06-10T06:48:51Z"
}
```
