/*
Package template generates json from multi jsons (json pool).

For example,
  giving the following jsons:
    json name: foo
    json data:
    {
      "a": 1
    }

    json name: bar
    json data:
    {
      "b": "apple"
    }

  giving template:
    {
      "w": 2019,
      "x": "*foo/a",
      "y": "*bar/b",
      "z": "\\*simple string"
    }

  result:
    {
      "w": 2019,
      "x": 1,
      "y": "apple",
      "z": "*simple string"
    }

Note: If the string in template prefixing with '*', the string will be
interpreted as 'named json pointer'. If a plain string in template prefix with
'*', it must be escaped, e.g., for string "\\*abc", after executing, the value
will be "*abc". For 'named json pointer', it is defined as following:

  named-json-pointer = json-name + json-pointer

If json name is 'foo', json-pointer is '/a', the named-json-pointer will be
'foo/a'.

To match each item in an array, we can use '$' as wildcard in
named-json-pointer. E.g.,

  giving the following json:
    json name: moo:
    json value:
    {
      "scores": [{"a":1}, {"a":2}, {"a":3}]
    }

  giving template:
    {
      "z": "*moo/scores/$/a"
    }

  result:
    {
      "z": [1, 2, 3]
    }

For more complicated example, please refer testcase/cases.json. In the json,
each case has field 'schema.const' which defined the json executed by the
template.
*/
package template
