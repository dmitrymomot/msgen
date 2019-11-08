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

## Usage

Let's create simple HTTP-RPC microservice with one method and exposed port via kubernetes load balancer service
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
Check it in terminal
```
curl --request POST \
  --url http://localhost:8888/twirp/testhttpsrv.Service/TestCall \
  --header 'content-type: application/json' \
  --data '{
	"str": "test string sample"
}'
```

## License

This library is licensed under the [Apache 2.0 License](https://github.com/dmitrymomot/msgen/blob/master/LICENSE).