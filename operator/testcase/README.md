# Testcase JSON Format
```json
[{
    "name": "The name of testcase",
    "enable": "Type is bool, enable the testcase or not",
    "cases": [{
        "data": "Test data",
        "expected": "Optional, expected result"
        "error": "Optional, type is bool, expect that the testcase will fail if true"
    }]
}]
```
