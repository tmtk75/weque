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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/spf13/cobra"
	"github.com/tmtk75/weque/repository"
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
	Short: "Print or send notification payload for testing",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		b, err := ioutil.ReadFile("./github/payload.json")
		if err != nil {
			log.Fatalf("%v", err)
		}
		//fmt.Printf("%v", string(b))
		var wh repository.Webhook
		err = json.Unmarshal(b, &wh)
		if err != nil {
			log.Fatalf("%v", err)
		}
		iwh, err := slack.NewIncomingWebhookRepository(&wh, &github.Github{}, nil)

		s, err := json.Marshal(iwh)
		if err != nil {
			log.Fatalf("%v", err)
		}
		fmt.Printf("%v", string(s))
	},
}

func init() {
	RootCmd.AddCommand(slackCmd)

	slackCmd.AddCommand(slackGithubCmd)
}
