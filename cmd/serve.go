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

	pflags.StringP(server.KeyPort, "p", ":9981", "port to listen")
	viper.BindPFlag(server.KeyPort, pflags.Lookup(server.KeyPort))

	pflags.Bool(server.KeyTLSEnabled, false, "enable TLS")
	viper.BindPFlag(server.KeyTLSEnabled, pflags.Lookup(server.KeyTLSEnabled))

	pflags.String(server.KeyTLSPort, ":https", "port in TLS to listen")
	viper.BindPFlag(server.KeyTLSPort, pflags.Lookup(server.KeyTLSPort))

	pflags.Bool(server.KeyACMEEnabled, false, "enable ACME to use Let's Encrypt")
	viper.BindPFlag(server.KeyACMEEnabled, pflags.Lookup(server.KeyACMEEnabled))

	pflags.String(server.KeyACMEPort, ":http", "port to listen for ACME challenge")
	viper.BindPFlag(server.KeyACMEPort, pflags.Lookup(server.KeyACMEPort))

	pflags.String(weque.KeyPrefix, "", "prefix for environment variable")
	viper.BindPFlag(weque.KeyPrefix, pflags.Lookup(weque.KeyPrefix))

	pflags.String(repositoryworker.KeyHandlerScriptRepository, "./handlers/repository", "handler script for repository")
	viper.BindPFlag(repositoryworker.KeyHandlerScriptRepository, pflags.Lookup(repositoryworker.KeyHandlerScriptRepository))

	pflags.String(registryworker.KeyHandlerScriptRegistry, "./handlers/registry", "handler script for registry")
	viper.BindPFlag(registryworker.KeyHandlerScriptRegistry, pflags.Lookup(registryworker.KeyHandlerScriptRegistry))

	pflags.Bool(weque.KeyInsecureMode, false, "skip steps to verify for development")
	viper.BindPFlag(weque.KeyInsecureMode, pflags.Lookup(weque.KeyInsecureMode))
	viper.BindEnv(weque.KeyInsecureMode, "INSECURE")
}
