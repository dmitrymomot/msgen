/*
Copyright © 2019 Dmytro Momot

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"
	"os"

	"github.com/dmitrymomot/msgen/src"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "msgen",
	Short:   "Generate a new microservice from template",
	Example: "msgen --grpc --grpc_port=9876 --pub --path=srv/users github.com/your-username/users-srv",
	Run:     src.Generate,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a service name")
		}
		if !cmd.Flags().Changed("grpc") && !cmd.Flags().Changed("twirp") && !cmd.Flags().Changed("pub") && !cmd.Flags().Changed("sub") {
			return errors.New("you should set at least one of server type: grpc, twirp, pub or sub")
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().StringP("namespace", "n", "default", "service namespace")
	rootCmd.Flags().String("path", "", "create in directory")

	rootCmd.Flags().Bool("linkerd", false, "enable linkerd support")

	rootCmd.Flags().Bool("grpc", false, "set up gRPC server")
	rootCmd.Flags().Bool("grpc_lb", false, "gRPC load balancer with external port")
	rootCmd.Flags().Bool("grpc_srv", false, "gRPC service to use only in cluster")
	rootCmd.Flags().Int("grpc_port", 9200, "grpc port")
	rootCmd.Flags().StringSlice("grpc_clients", nil, "add grpc clients connections (e.g.: --grpc_clients=users-srv:9200,order:9200)")
	rootCmd.Flags().Bool("grpc_client_helper", false, "create gRPC client helper")

	rootCmd.Flags().StringSlice("rpc_methods", nil, "generate rpc methods (e.g.: --rpc_methods=list_users,get_user_by_id)")

	rootCmd.Flags().Bool("twirp", false, "set up twirp server (exposed via http_port)")

	rootCmd.Flags().Bool("http", false, "set up http server")
	rootCmd.Flags().Int("http_port", 8888, "http port")
	rootCmd.Flags().Bool("http_lb", false, "HTTP k8s load balancer with external port")
	rootCmd.Flags().Bool("http_srv", false, "HTTP k8s service to use only in cluster")

	rootCmd.Flags().Bool("pub", false, "set up NATS publisher")
	rootCmd.Flags().Bool("sub", false, "set up NATS subscriber")

	rootCmd.Flags().Bool("redis_pool", false, "set up redis connections pool")
	rootCmd.Flags().String("redis_host", "", "redis connection path")

	rootCmd.Flags().StringSlice("jobs", nil, "set up jobs queue based on redis (e.g.: --jobs=send_email,send_notification)")

	rootCmd.Flags().String("db_host", "postgres", "database host")
	rootCmd.Flags().Int("db_port", 5432, "database port")
	rootCmd.Flags().String("db_name", "", "database name")
	rootCmd.Flags().String("db_user", "", "database user")
	rootCmd.Flags().String("db_password", "", "database password")

	rootCmd.Flags().Bool("k8s", true, "generate kubernetes deployment config")
	rootCmd.Flags().Bool("tls", false, "add tls certificate loader")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match
}
