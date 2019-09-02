# Introduction
This project implemented a set of tools to create, manipulate and validate JSONs. It was created when I was developing a project (closed source, inspired by [Blueprints of UnrealEngine](https://docs.unrealengine.com/en-US/Engine/Blueprints/index.html)) which uses configuration files to define RESTful APIs. The project splited API handling to multiple stages (e.g., processing request stage, MySQL/MongoDB interact stage, Redis interact stage, 3rd-party API call stage, response stage, etc), all these stages share one JSON pool, and each stage can create the needed JSON as input from the JSON pool (by using template and operator), also, each stage can add its product JSON to the JSON pool for other stages to use.

* [JSON template](./template) - Create JSON from multiple JSONs
* [JSON schema](./schema) - Validate JSON against schema
* [JSON operator](./operator) - Manipulate JSON


# Playground
1. [template](https://play.golang.org/p/3JcUvPEvur7)
2. [template & operator](https://play.golang.org/p/GWHASNc_BhN)
