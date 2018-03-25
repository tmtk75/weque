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

	"github.com/tmtk75/weque/consumer"
	"github.com/tmtk75/weque/registry"
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

	pflags.IntP("port", "p", 9981, "port to listen")
	viper.BindPFlag("port", pflags.Lookup("port"))

	pflags.String("prefix", "", "prefix for environment variable")
	viper.BindPFlag("prefix", pflags.Lookup("prefix"))

	pflags.String("repository-handler", "./handlers/repository", "handler script for repository")
	viper.BindPFlag(consumer.KeyHandlerScriptRepository, pflags.Lookup("repository-handler"))

	pflags.String("registry-handler", "./handlers/registry", "handler script for registry")
	viper.BindPFlag(registry.KeyHandlerScript, pflags.Lookup("registry-handler"))
}