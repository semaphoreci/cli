package cmd

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/ghodss/yaml"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Get and set configuration options.",
	Long:  ``,
}

var configGetCmd = &cobra.Command{
	Use:   "get [NAME]",
	Short: "Display a configuration options",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		home, err := homedir.Dir()
		path := fmt.Sprintf("%s/.sem.yaml", home)
		data, err := ioutil.ReadFile(path)

		m := make(map[string]string)

		err = yaml.Unmarshal([]byte(data), &m)

		if err != nil {
			log.Fatalf("error: %v", err)
		}

		fmt.Println(m[name])
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set [NAME] [VALUE]",
	Short: "Set a configuration options",
	Long:  ``,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		value := args[1]

		home, err := homedir.Dir()
		path := fmt.Sprintf("%s/.sem.yaml", home)
		data, err := ioutil.ReadFile(path)

		m := make(map[string]string)

		err = yaml.Unmarshal([]byte(data), &m)

		if err != nil {
			fmt.Printf("sdasdsa")
		}

		m[name] = value

		d, err := yaml.Marshal(&m)

		if err != nil {
			log.Fatalf("error: %v", err)
		}

		err = ioutil.WriteFile(path, d, 0644)

		if err != nil {
			log.Fatalf("error: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
}
