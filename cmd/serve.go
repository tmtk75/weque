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

	// string options
	stropts := []struct {
		k    string
		v    string
		desc string
	}{
		{k: server.KeyPort, v: ":9981", desc: "port to listen"},
		{k: server.KeyTLSPort, v: ":https", desc: "port in TLS to listen"},
		{k: server.KeyTLSCertFile, v: "", desc: "cert file to listen in TLS"},
		{k: server.KeyTLSKeyFile, v: "", desc: "keyfile to listen in TLS"},
		{k: server.KeyACMEPort, v: ":http", desc: "port to listen for ACME challenge"},
		{k: weque.KeyPrefix, v: "", desc: "prefix for environment variable"},
		{k: repositoryworker.KeyHandlerScriptRepository, v: "./handlers/repository", desc: "handler script for repository"},
		{k: registryworker.KeyHandlerScriptRegistry, v: "./handlers/registry", desc: "handler script for registry"},
	}

	// boolean options
	boolopts := []struct {
		k    string
		v    bool
		desc string
	}{
		{k: server.KeyTLSEnabled, v: false, desc: "enable TLS"},
		{k: server.KeyACMEEnabled, v: false, desc: "enable ACME to use Let's Encrypt"},
		{k: weque.KeyInsecureMode, v: false, desc: "skip steps to verify for development"},
	}

	for _, e := range stropts {
		pflags.String(e.k, e.v, e.desc)
		viper.BindPFlag(e.k, pflags.Lookup(e.k))
	}

	for _, e := range boolopts {
		pflags.Bool(e.k, e.v, e.desc)
		viper.BindPFlag(e.k, pflags.Lookup(e.k))
	}

	viper.BindEnv(weque.KeyInsecureMode, "INSECURE")
}
