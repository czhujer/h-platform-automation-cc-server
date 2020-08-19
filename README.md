# h-platform-automation-cc-server
Command and Control server for whole platform

## main docs 
- https://github.com/czhujer/h-platform-automation-docs/blob/master/README.md

## development

### docs
- https://medium.com/the-andela-way/build-a-restful-json-api-with-golang-85a83420c9da
- https://medium.com/@gsisimogang/instrumenting-golang-server-in-5-min-c1c32489add3

## commands
```
curl  -d '{"key1":"value1", "key2":"value2"}' -H "Content-Type: application/json" -X POST 127.0.0.1:8080/calculoid/webhook -i
```
```
curl  -d @test-data.json -H "Content-Type: application/json" -X POST 127.0.0.1:8080/calculoid/webhook -i
```
