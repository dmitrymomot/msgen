# msgen

The simple tool to generate basic boilerplate of microservice. Not some special, just for my own needs.

Installation
```
go get -u github.com/dmitrymomot/msgen
```

Available options
```
msgen -h
```

## Features

* gRPC server & client
* [Twirp](https://github.com/twitchtv/twirp) server
* Pub/sup via [NATS](https://nats.io)
* Jobs queue via redis

## Maybe better alternatives for you

* https://github.com/lileio/lile (Easily generate gRPC services in Go)
* https://github.com/izumin5210/grapi (A surprisingly easy API server and generator in gRPC and Go)
* https://micro.mu/docs/new.html (a part of the `micro` infrastructure)
* https://github.com/fiorix/protoc-gen-cobra (Cobra command line tool generator for gRPC clients)

## Prerequisites

To use the tool you should have installed [brotobuf](https://developers.google.com/protocol-buffers/docs/gotutorial) and [grpc-go](https://github.com/grpc/grpc-go)

```
go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
go get -u google.golang.org/grpc
```
Protobuf is also available in MacOS through Homebrew:
```
brew install protobuf
```

### Optional

[Twirp](https://twitchtv.github.io/twirp/docs/install.html) - is a simple RPC framework built on protobuf
```
go get -u github.com/twitchtv/twirp/protoc-gen-twirp
```
Generate twirp files
```
make tproto
```

[Protobuf validator](https://github.com/envoyproxy/protoc-gen-validate) by Envoyproxy*
```
# fetches this repo into $GOPATH
go get -d github.com/envoyproxy/protoc-gen-validate

# installs PGV into $GOPATH/bin
cd $GOPATH/src/github.com/envoyproxy/protoc-gen-validate && make build
```
Generate validator for your proto
```
make vproto
```

## Usage

Let's create simple HTTP-RPC microservice with one method and exposed port via kubernetes load balancer service.
To move ahead through the usage example you should install Twirp from [optional dependencies](#optional) section.
> Run `msgen --help` to get more details about available options
```
msgen --twirp --rpc_methods=test_call --http_lb test-http-srv
```
The command above generates next files
```
├── Dockerfile
├── Makefile
├── README.md
├── go.mod
├── k8s.yml
├── logger
│   └── logger.go
├── main.go
├── pb
│   └── testhttpsrv
│       └── service.proto
└── service
    ├── service.go
    └── test_call.go
```
Move into test-http-srv and run service
```
cd test-http-srv && make build docker && docker run -d  -p 8888:8888 test-http-srv:latest
```
or if you have installed kubernetes
```
cd test-http-srv && make build docker deploy
```
Check it in terminal
```
curl --request POST \
  --url http://localhost:8888/twirp/testhttpsrv.Service/TestCall \
  --header 'content-type: application/json' \
  --data '{"str": "test string sample"}'
```

## License

This library is licensed under the [Apache 2.0 License](https://github.com/dmitrymomot/msgen/blob/master/LICENSE).