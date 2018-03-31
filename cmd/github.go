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
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/alecthomas/template"
	"github.com/hashicorp/hcl"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tmtk75/weque/github"
)

// githubCmd represents the github command
var githubCmd = &cobra.Command{
	Use:   "github",
	Short: "Commands for GitHub",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	// github
	RootCmd.AddCommand(githubCmd)

	// github list
	githubCmd.AddCommand(listCmd)

	// github create
	githubCmd.AddCommand(createCmd)
	createCmd.PersistentFlags().String("name", "web", "webhook name")

	// github tf
	githubCmd.AddCommand(tfCmd)

	// github tfgen
	githubCmd.AddCommand(tfgenCmd)
	tfgenCmd.PersistentFlags().StringP("content-type", "c", "application/json", "content-type")
	viper.BindPFlag("content_type", tfgenCmd.PersistentFlags().Lookup("content-type"))
}

var listCmd = &cobra.Command{
	Use: "list <onwer/repo>",
	Long: `To give your token,

    export GITHUB_TOKEN=...`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		github.List(args[0])
	},
}

var createCmd = &cobra.Command{
	Use:  "create <repo> <url> <secret>",
	Long: `Create a new webhook in the github repository.`,
	Args: cobra.RangeArgs(3, 4), // REPO URL SECRET
	Run: func(cmd *cobra.Command, args []string) {
		var (
			repo   = args[0]
			url    = args[1]
			secret = args[2]
			name   = "web"
		)
		if len(args) == 4 {
			name = args[3]
		}
		github.Create(repo, url, secret, name)
	},
}

var tfCmd = &cobra.Command{
	Use: "tf",
	Run: func(cmd *cobra.Command, args []string) {
		a, _ := ioutil.ReadAll(os.Stdin)
		var out interface{}
		err := hcl.Decode(&out, string(a))
		if err != nil {
			log.Fatalf("%v", err)
		}
		fmt.Printf("%v\n", out)
	},
}

var tfgenCmd = &cobra.Command{
	Use:  "tfgen",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		t := template.New("")
		t, _ = t.Parse(tmpl)

		ctype := cmd.Flag("content-type")
		fmt.Printf("ctype: %v\n", ctype.Value)

		t.Execute(os.Stdout, map[string]string{
			"Name":        args[0],
			"URL":         args[1],
			"ContentType": viper.GetString("content_type"),
		})
	},
}

var tmpl = `#
variable "repository_name" { default = "{{ .Name }}" }

resource "github_repository_webhook" "foo" {
  repository = "${var.repository_name}"

  name   = "web"
  active = true
  events = ["push"]

  configuration {
    url          = "{{ .URL }}"
    secret       = "${var.webhook_secret}"
    content_type = "{{ .ContentType }}"
    insecure_ssl = false
  }
}
`
