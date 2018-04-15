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
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tmtk75/weque/repository/bitbucket"
	"github.com/tmtk75/weque/repository/github"
	"github.com/tmtk75/weque/slack"
)

// slackCmd represents the slack command
var slackCmd = &cobra.Command{
	Use:   "slack",
	Short: "A brief description of your command",
	Long:  `To give your incoming url, use SLACK_URL`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var slackGithubCmd = &cobra.Command{
	Use:   "github [flags]",
	Short: "Print slack payload for GitHub",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		r, _ := http.NewRequest("POST", "http://example.com", nil)
		r.Header.Add("Content-type", "application/json")
		r.Header.Add("X-Github-Delivery", "")
		e := &github.Github{}
		slack.PrintIncomingWebhookRepository(r, "./github/payload.json", e, e)
	},
}

var slackBitbucketCmd = &cobra.Command{
	Use:   "bitbucket [flags]",
	Short: "Print slack payload for Bitbucket",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		r, _ := http.NewRequest("POST", "http://example.com", nil)
		r.Header.Add("X-Event-Key", "a")
		r.Header.Add("X-Request-UUID", "b")
		e := &bitbucket.Bitbucket{}
		slack.PrintIncomingWebhookRepository(r, "./bitbucket/payload.json", e, e)
	},
}

var slackRegistryCmd = &cobra.Command{
	Use:   "registry [flags]",
	Short: "Print slack payload for docker registry",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		slack.PrintIncomingWebhookRegistry("./registry/payload.json")
	},
}

func init() {
	RootCmd.AddCommand(slackCmd)
	pflags := slackCmd.PersistentFlags()
	pflags.String("channel-name", "#api-test", "Slack channel name")
	viper.BindPFlag(slack.KeySlackChannelName, pflags.Lookup("channel-name"))

	//
	slackCmd.AddCommand(slackGithubCmd)
	slackCmd.AddCommand(slackBitbucketCmd)
	slackCmd.AddCommand(slackRegistryCmd)
}
