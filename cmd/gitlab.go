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
	"github.com/spf13/cobra"

	"github.com/tmtk75/weque/gitlab"
)

// gitlabCmd represents the gitlab command
var gitlabCmd = &cobra.Command{
	Use:   "gitlab",
	Short: "Commands for gitlab",
	Long: `
  export GITLAB_PRIVATE_TOKEN=...
  Issue it via https://gitlab.com/profile/personal_access_tokens
`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List hooks for the given project",
	Long:  `list tmtk75/foobar`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		gitlab.List(args[0])
	},
}

var createCmd = &cobra.Command{
	Use:   "create [flags] <project> <URI> <secret>",
	Short: "Create a new hook for the given project",
	Long:  `create tmtk75/foobar https://example.com`,
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		gitlab.Create(args[0], args[1], args[2])
	},
}

func init() {
	RootCmd.AddCommand(gitlabCmd)

	gitlabCmd.AddCommand(listCmd)
	gitlabCmd.AddCommand(createCmd)
}
