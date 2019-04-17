package cmd

import (
	"fmt"
	"regexp"
	"strings"

	client "github.com/semaphoreci/cli/api/client"
	models "github.com/semaphoreci/cli/api/models"
	"github.com/semaphoreci/cli/cmd/utils"
	"github.com/spf13/cobra"
)

func NewCreateSecretCmd() *cobra.Command {
	cmd := &cobra.Command{}

	cmd.Use = "secret [NAME]"
	cmd.Short = "Create a secret."
	cmd.Long = ``
	cmd.Aliases = []string{"secrets"}
	cmd.Args = cobra.ExactArgs(1)

	cmd.Flags().StringArrayP(
		"file",
		"f",
		[]string{},
		"File mapping <local-path>:<mount-path>, used to create a secret with file",
	)

	cmd.Flags().StringArrayP(
		"env",
		"e",
		[]string{},
		"Environment Variables",
	)

	cmd.Run = func(cmd *cobra.Command, args []string) {
		name := args[0]

		fileFlags, err := cmd.Flags().GetStringArray("file")
		utils.Check(err)

		var files []models.SecretV1BetaFile
		for _, fileFlag := range fileFlags {
			matchFormat, err := regexp.MatchString(`^[^: ]+:[^: ]+$`, fileFlag)
			utils.Check(err)

			if !matchFormat {
				utils.Fail("The format of --file flag must be: <local-path>:<semaphore-path>")
			}

			flagPaths := strings.Split(fileFlag, ":")

			file := models.SecretV1BetaFile{}
			file.Path = flagPaths[1]
			file.Content = encodeFromFileAt(flagPaths[0])
			files = append(files, file)
		}

		envFlags, err := cmd.Flags().GetStringArray("env")
		utils.Check(err)

		var envVars []models.SecretV1BetaEnvVar
		for _, envFlag := range envFlags {
			matchFormat, err := regexp.MatchString(`^.+=.+$`, envFlag)
			utils.Check(err)

			if !matchFormat {
				utils.Fail("The format of -e flag must be: <NAME>=<VALUE>")
			}

			parts := strings.SplitN(envFlag, "=", 2)

			envVars = append(envVars, models.SecretV1BetaEnvVar{
				Name:  parts[0],
				Value: parts[1],
			})
		}

		secret := models.NewSecretV1Beta(name, envVars, files)

		c := client.NewSecretV1BetaApi()

		_, err = c.CreateSecret(&secret)

		utils.Check(err)

		fmt.Printf("Secret '%s' created.\n", secret.Metadata.Name)
	}

	return cmd
}
