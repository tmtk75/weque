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
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/tmtk75/weque"
	gh "github.com/tmtk75/weque/github"
	"github.com/tmtk75/weque/repository/github"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run github, bitbucket and registry",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

//go:generate go-assets-builder -p github -o ../github/payload.go ../github/payload.json

var runGithubCmd = &cobra.Command{
	Use:   "github [flags] <command> [args...]",
	Args:  cobra.MinimumNArgs(1),
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		r, _ := http.NewRequest("POST", "http://example.com", nil)
		r.Header.Add("Content-type", "application/json")
		r.Header.Add("X-Github-Delivery", "delivery-id")
		e := &github.Github{}

		a, _ := gh.Assets.Open("/github/payload.json")
		b, _ := ioutil.ReadAll(a)
		wh, _ := e.Unmarshal(r, b)
		weque.Stdout = os.Stdout
		weque.Stderr = os.Stderr
		weque.Run(wh.Env(), ".", args[0], args[1:len(args)]...)
	},
}

func init() {
	RootCmd.AddCommand(runCmd)
	runCmd.AddCommand(runGithubCmd)
}
