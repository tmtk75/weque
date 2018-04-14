// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tmtk75/weque"
	"github.com/tmtk75/weque/registry/worker"
	"github.com/tmtk75/weque/repository/worker"
	"github.com/tmtk75/weque/server"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start a server",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if err := server.Validate(); err != nil {
			log.Printf("[error] %v", err)
			return
		}
		s := server.New()
		s.Start()
	},
}

func init() {
	RootCmd.AddCommand(serverCmd)

	pflags := serverCmd.PersistentFlags()

	pflags.StringP("port", "p", ":9981", "port to listen")
	viper.BindPFlag(server.KeyPort, pflags.Lookup("port"))

	pflags.String("tls.port", ":https", "port in TLS to listen")
	viper.BindPFlag(server.KeyTLSPort, pflags.Lookup("tls.port"))

	pflags.Bool("acme.enabled", false, "Enable Let's Encrypt")
	viper.BindPFlag(server.KeyACMEEnabled, pflags.Lookup("acme.enabled"))

	pflags.String("acme.challenge-port", ":http", "port to listen for ACME challenge")
	viper.BindPFlag(server.KeyACMEChallengePort, pflags.Lookup("acme.challenge-port"))

	pflags.String("prefix", "", "prefix for environment variable")
	viper.BindPFlag(weque.KeyPrefix, pflags.Lookup("prefix"))

	pflags.String("repository-handler", "./handlers/repository", "handler script for repository")
	viper.BindPFlag(repositoryworker.KeyHandlerScriptRepository, pflags.Lookup("repository-handler"))

	pflags.String("registry-handler", "./handlers/registry", "handler script for registry")
	viper.BindPFlag(registryworker.KeyHandlerScriptRegistry, pflags.Lookup("registry-handler"))

	pflags.Bool("insecure", false, "handler script for registry")
	viper.BindPFlag(weque.KeyInsecureMode, pflags.Lookup("insecure"))
	viper.BindEnv(weque.KeyInsecureMode, "INSECURE")
}
