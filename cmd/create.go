package cmd

import (
	"fmt"
  "io/ioutil"
  "os"
  "net/http"
  "bytes"

	"github.com/spf13/cobra"
  "github.com/ghodss/yaml"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a resource from a file.",
	Long: ``,

	Run: func(cmd *cobra.Command, args []string) {
    RunCreate(cmd, args)

	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.PersistentFlags().StringP("file", "f", "", "Filename, directory, or URL to files to use to create the resource")
}

func RunCreate(cmd *cobra.Command, args []string) {
    path, err := cmd.Flags().GetString("file")

    check(err, "Path not provided")

    data, err := ioutil.ReadFile(path)

    check(err, "Failed to read from resource file.")

    resource, err := parse(data)

    check(err, "Failed to parse resource file.")

    apiVersion := resource["apiVersion"].(string)
    kind := resource["kind"].(string)

    upload(apiVersion, kind, data)
}

func parse(data []byte) (map[string]interface{}, error) {
  m := make(map[string]interface{})

  // fmt.Print(string(data))

  err := yaml.Unmarshal(data, &m)

  return m, err
}

func upload(version string, kind string, data []byte) {
  // fmt.Printf("apiVersion: %s\n", version)
  // fmt.Printf("kind: %s\n", kind)

  var path string

  switch kind {
    case "Secret":
      path = fmt.Sprintf("/api/%s/secrets", version)

    default:
      panic("Unsuported kind")
  }

  url := fmt.Sprintf("http://renderedtext.semaphoreci.com%s", path)

  // fmt.Printf("Path: %s\n", url)

  j, err := yaml.YAMLToJSON(data)

  // fmt.Printf("Content : %s\n", j)

  req, err := http.NewRequest("POST", url, bytes.NewBuffer(j))

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

func check(err error, message string) {
  if err != nil {
    fmt.Fprintf(os.Stderr, "%s\n", message)
    fmt.Fprintf(os.Stderr, "error: %v\n", err)

    os.Exit(1)
  }
}
