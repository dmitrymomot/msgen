syntax = "proto3";

package {{ package .ServiceName }};

// available plugins:
// import "google/protobuf/wrappers.proto";
// import "google/api/annotations.proto";
// import "validate/validate.proto";
// import "google/protobuf/any.proto";

service Service {
	{{if .RPCMethods}}{{range $method := .RPCMethods}}
	rpc {{title $method}} ({{title $method}}Req) returns ({{title $method}}Resp) {}{{end}}{{end}}
}

{{if .RPCMethods}}{{range $method := .RPCMethods}}
message {{title $method}}Req {
    string str = 1;
}

message {{title $method}}Resp {
	string str = 1;
}
{{end}}{{end}}