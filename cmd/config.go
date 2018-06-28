package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

		if viper.IsSet(name) {
			value := viper.GetString(name)

			fmt.Println(value)
		} else {
			fmt.Printf("configuration \"%s\" not found\n", name)

			os.Exit(1)
		}
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

		viper.Set(name, value)
		viper.WriteConfig()
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
}
