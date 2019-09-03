# Introduction
This project implemented a set of tools to create, manipulate and validate JSONs. It was created when I was developing a project (closed source, inspired by [Blueprints of UnrealEngine](https://docs.unrealengine.com/en-US/Engine/Blueprints/index.html)) which uses configuration files to define RESTful APIs. The project splited API handling to multiple stages (e.g., processing request stage, MySQL/MongoDB interact stage, Redis interact stage, 3rd-party API call stage, response stage, etc), all these stages share one JSON pool, and each stage can create the needed JSON as input from the JSON pool (by using template and operator), also, each stage can add its product JSON to the JSON pool for other stages to use.

* [JSON template](./template) - Create JSON from multiple JSONs
* [JSON schema](./schema) - Validate JSON against schema
* [JSON operator](./operator) - Manipulate JSON


## Example
Giving JSON 'foo':
```json
{
    "a1": 1,
    "a2": [1, 2, 3]
}
```

Giving JSON 'bar':
```json
{
    "b1": "hello",
    "b2": "world",
    "b3": ["a", "b"]
}
```

Giving JSON template:
```json
{
    "f1": ["*foo/a1", "*bar/b1", "*bar/b2"],
    "f2": [{"f21": "*foo/a2/$", "f22": "*bar/b3/$"}],
    "f3": "*bar/b3/1"
}
```

After executing the template, we can get the following JSON:
```json
{
    "f1": [1, "hello", "world"],
    "f2": [{"f21": 1, "f22": "a"}, {"f21": 2, "f22": "b"}, {"f21": 3}],
    "f3": "b"
}
```

## Playground
1. [template](https://play.golang.org/p/3JcUvPEvur7)
2. [template & operator](https://play.golang.org/p/GWHASNc_BhN)
