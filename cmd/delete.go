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
  "net/http"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a resource.",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
    RunDelete(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func RunDelete(cmd *cobra.Command, args []string) {
  var path string

  version := "v1alpha"

  kind := args[0]
  name := args[1]

  switch kind {
    case "secret":
      path = fmt.Sprintf("/api/%s/secrets/%s", version, name)

    case "secrets":
      path = fmt.Sprintf("/api/%s/secrets/%s", version, name)

    default:
      panic("Unsuported kind")
  }

  url := fmt.Sprintf("http://renderedtext.semaphoreci.com%s", path)

  req, err := http.NewRequest("DELETE", url, nil)

  req.Header.Set("Content-Type", "application/json")
  req.Header.Set("X-Semaphore-Req-ID", "111")
  req.Header.Set("X-Semaphore-User-ID", "111")
  req.Header.Set("Authorization", "Token C4V6j96w7D5YHqWJGHxz")

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
      panic(err)
  }
  defer resp.Body.Close()

  // fmt.Println("response Status:", resp.Status)
  // fmt.Println("response Headers:", resp.Header)
  body, _ := ioutil.ReadAll(resp.Body)
  fmt.Println(string(body))
}
