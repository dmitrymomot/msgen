syntax = "proto3";

package {{ grpcpackage .ServiceName }};
option go_package="{{ package .ServiceName }}";

// available plugins:
// import "google/protobuf/wrappers.proto";
// import "validate/validate.proto";
// import "google/protobuf/any.proto";
{{if not .GrpcGateway}}// {{end}}import "google/api/annotations.proto";

service Service {
	{{if .RPCMethods}}{{range $method := .RPCMethods}}
	rpc {{title $method}} ({{title $method}}Req) returns ({{title $method}}Resp) { {{if $.GrpcGateway}}
		option (google.api.http) = {
			post: "/{{if not (eq $.Version "")}}{{$.Version}}/{{end}}{{ package $.ServiceName }}/{{url $method}}"
			body: "*"
		};
	{{end}}}{{end}}{{end}}
}

{{if .RPCMethods}}{{range $method := .RPCMethods}}
message {{title $method}}Req {
    string str = 1;
}

message {{title $method}}Resp {
	string str = 1;
}
{{end}}{{end}}