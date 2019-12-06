package src

import (
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type (
	options struct {
		ServiceName string
		ServicePath string
		Namespace   string
		Version     string
		Path        string
		AbsPath     string

		Linkerd bool

		Grpc        bool
		GrpcLB      bool
		GrpcPort    int
		GrpcClients []string
		GrpcSrv     bool
		GrpcGateway bool

		Twirp bool

		RPC        bool
		RPCMethods []string

		HTTP          bool
		HTTPPort      int
		HTTPEndpoints []string
		HTTPLB        bool
		HTTPSrv       bool

		RedisPool bool
		RedisHost string
		Jobs      []string

		DB *dbOptions

		ClientHelper bool
		K8s          bool

		TLS bool

		Nats bool
		Pub  bool
		Sub  bool

		CustomOptions map[string]interface{}
	}

	dbOptions struct {
		Host     string
		Port     int
		Name     string
		User     string
		Password string
	}
)

func parseOptions(cmd *cobra.Command, args []string) options {
	opt := options{}

	opt.ServicePath = cmd.Flags().Arg(0)
	opt.ServiceName = filepath.Base(opt.ServicePath)
	opt.Namespace, _ = cmd.Flags().GetString("namespace")
	opt.Version, _ = cmd.Flags().GetString("version")

	opt.Path, _ = cmd.Flags().GetString("path")
	opt.Path = strings.Trim(opt.Path, "/")
	if opt.Path == "" || opt.Path == "." {
		opt.Path = opt.ServiceName
	}
	opt.AbsPath, _ = filepath.Abs(opt.Path)

	opt.Linkerd, _ = cmd.Flags().GetBool("linkerd")

	opt.Grpc, _ = cmd.Flags().GetBool("grpc")
	opt.GrpcLB, _ = cmd.Flags().GetBool("grpc_lb")
	opt.GrpcPort, _ = cmd.Flags().GetInt("grpc_port")
	opt.GrpcSrv, _ = cmd.Flags().GetBool("grpc_srv")
	opt.GrpcGateway, _ = cmd.Flags().GetBool("grpc_gateway")
	opt.GrpcClients, _ = cmd.Flags().GetStringSlice("grpc_clients")

	if opt.Grpc && !opt.GrpcLB && !opt.GrpcSrv {
		opt.GrpcSrv = true
	}

	opt.Twirp, _ = cmd.Flags().GetBool("twirp")

	opt.RPC = opt.Twirp || opt.Grpc
	opt.RPCMethods, _ = cmd.Flags().GetStringSlice("rpc_methods")

	opt.HTTP, _ = cmd.Flags().GetBool("http")
	opt.HTTPPort, _ = cmd.Flags().GetInt("http_port")
	opt.HTTPEndpoints, _ = cmd.Flags().GetStringSlice("http_endpoints")
	opt.HTTPLB, _ = cmd.Flags().GetBool("http_lb")
	opt.HTTPSrv, _ = cmd.Flags().GetBool("http_srv")
	opt.HTTP = opt.Twirp || opt.HTTP

	if opt.HTTP && !opt.HTTPLB && !opt.HTTPSrv {
		opt.HTTPSrv = true
	}

	opt.RedisPool, _ = cmd.Flags().GetBool("redis_pool")
	opt.RedisHost, _ = cmd.Flags().GetString("redis_host")
	opt.Jobs, _ = cmd.Flags().GetStringSlice("jobs")
	opt.RedisPool = opt.RedisPool || len(opt.Jobs) > 0

	dbopt := dbOptions{}
	dbopt.Host, _ = cmd.Flags().GetString("db_host")
	dbopt.Port, _ = cmd.Flags().GetInt("db_port")
	dbopt.Name, _ = cmd.Flags().GetString("db_name")
	dbopt.User, _ = cmd.Flags().GetString("db_user")
	dbopt.Password, _ = cmd.Flags().GetString("db_password")

	if dbopt.Name != "" && dbopt.Host != "" && dbopt.Port > 0 {
		opt.DB = &dbopt
	}

	opt.ClientHelper, _ = cmd.Flags().GetBool("client_helper")
	opt.K8s, _ = cmd.Flags().GetBool("k8s")
	opt.K8s = opt.K8s || opt.HTTPLB || opt.HTTPSrv || opt.GrpcSrv || opt.GrpcLB

	opt.TLS, _ = cmd.Flags().GetBool("tls")

	opt.Pub, _ = cmd.Flags().GetBool("pub")
	opt.Sub, _ = cmd.Flags().GetBool("sub")
	opt.Nats = opt.Pub || opt.Sub

	return opt
}
