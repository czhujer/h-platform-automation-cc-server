# h-platform-automation-cc-server
Command and Control server for whole platform

## main docs 
- https://github.com/czhujer/h-platform-automation-docs/blob/master/README.md

## development

### docs
- https://medium.com/the-andela-way/build-a-restful-json-api-with-golang-85a83420c9da
- https://medium.com/@gsisimogang/instrumenting-golang-server-in-5-min-c1c32489add3
- https://www.alexedwards.net/blog/making-and-using-middleware
- https://gowebexamples.com/advanced-middleware/

#### metrics
- https://www.robustperception.io/prometheus-middleware-for-gorilla-mux

### tracing
- https://opentracing.io/guides/golang/quick-start/
- https://medium.com/opentracing/tracing-http-request-latency-in-go-with-opentracing-7cc1282a100a
- https://github.com/opentracing-contrib/go-gorilla
- https://github.com/opentracing-contrib/go-gorilla/blob/master/gorilla/example_test.go
- https://medium.com/@carlosedp/instrumenting-go-for-tracing-c5bdabe1fc81

## commands
```
curl  -d '{"key1":"value1", "key2":"value2"}' -H "Content-Type: application/json" -X POST 127.0.0.1:8080/calculoid/webhook -i
```
```
curl  -d @test-data.json -H "Content-Type: application/json" -X POST 127.0.0.1:8080/calculoid/webhook -i
```
