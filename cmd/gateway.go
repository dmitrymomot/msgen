/*
Copyright © 2019 Dmytro Momot <mail@dmomot.com>

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
	"github.com/dmitrymomot/msgen/src"
	"github.com/spf13/cobra"
)

// gatewayCmd represents the gateway command
var gatewayCmd = &cobra.Command{
	Use:     "gateway",
	Aliases: []string{"gw", "api"},
	Short:   "Generate a new grpc gateway",
	Run:     src.GenerateGateway,
}

func init() {
	rootCmd.AddCommand(gatewayCmd)

	gatewayCmd.Flags().StringSlice("services", nil, "add grpc services list (e.g.: --services=users-srv:9200,order:9200)")
}
