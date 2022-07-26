# OpenAPI Specifications

All the products generated with [go-swagger](https://github.com/go-swagger/go-swagger) have to be written in [OpenAPI 2.0](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md)

## Adding a new specification

Specifications should be added under following convention.

For each service related to a product we should create a `${PRODUCT}` directory which should contain a `${SERVICE}.yaml` specification file, for nested or composed specifications subdirectories may be created as well.

The final repository structure should look like this:

```
.
+-- _foo
|   +-- _bar
|   |   +-- baz.yaml
|   |   ...
|   ...
```
