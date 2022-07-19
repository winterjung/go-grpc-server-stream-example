# go-grpc-server-stream-example

```shell
$ curl -sSL "https://github.com/bufbuild/buf/releases/download/v1.6.0/buf-$(uname -s)-$(uname -m)" -o ./buf
$ chmod +x ./buf
$ go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc
$ ./buf generate
$ go run -v -race cmd/main.go
$ curl -v http://localhost:8081/images/casper.png
> GET /images/casper.png HTTP/1.1
> Host: localhost:8081
> User-Agent: curl/7.84.0
> Accept: */*
>

< HTTP/1.1 200 OK
< Content-Type: application/json
< Grpc-Metadata-Content-Type: application/grpc
< Date: Tue, 19 Jul 2022 02:15:06 GMT
< Transfer-Encoding: chunked
{"result":{"data":"iVBORw0KGg..."}}
{"result":{"data":"ZwSGWU4BjR..."}}
{"result":{"data":"VpWuaZD81J..."}}
{"result":{"data":"JeYVF8Osyf..."}}
{"result":{"data":"S9eO92LRDW..."}}
{"result":{"data":"PBDt/G58ih..."}}
{"result":{"data":"sU39sPHAxS..."}}
{"result":{"data":"98jrORYdhn..."}}
{"result":{"data":"O4N4jJ+XCY..."}}
```

## refs

- https://ops.tips/blog/sending-files-via-grpc/
- https://dev.to/techschoolguru/upload-file-in-chunks-with-client-streaming-grpc-golang-4loc
- https://dev.to/techschoolguru/implement-server-streaming-grpc-in-go-2d0p
