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
	"github.com/tmtk75/weque/bitbucket"
)

// bitbucketCmd represents the bitbucket command
var bitbucketCmd = &cobra.Command{
	Use:   "bitbucket",
	Short: "Commands for Bitbucket",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	RootCmd.AddCommand(bitbucketCmd)

	bitbucketCmd.AddCommand(bitbucketListCmd)
	bitbucketCmd.AddCommand(bitbucketCreateCmd)
}

var bitbucketListCmd = &cobra.Command{
	Use:  "list [flags] <owner/repo>",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		bitbucket.List(args[0])
	},
}

var bitbucketCreateCmd = &cobra.Command{
	Use:  "create",
	Long: `Create a new webhook in the bitbucket repository.`,
	Args: cobra.ExactArgs(3), // REPO URL SECRET
	Run: func(cmd *cobra.Command, args []string) {
		var (
			repo   = args[0]
			url    = args[1]
			secret = args[2]
		)
		bitbucket.Create(repo, url, secret)
	},
}
