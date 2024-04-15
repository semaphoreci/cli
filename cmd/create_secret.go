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
	cmd.Short = "Create a secret. If project is specified (via -p or -i), create a project secret."
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

	cmd.Flags().StringP("project-name", "p", "",
		"project name; if specified will edit project level secret, otherwise organization secret")
	cmd.Flags().StringP("project-id", "i", "",
		"project id; if specified will edit project level secret, otherwise organization secret")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		projectID := GetProjectID(cmd)

		name := args[0]

		fileFlags, err := cmd.Flags().GetStringArray("file")
		utils.Check(err)
		envFlags, err := cmd.Flags().GetStringArray("env")
		utils.Check(err)

		if projectID == "" {
			files := parseSecretFiles(fileFlags)
			envVars := parseSecretEnvVars(envFlags)

			secret := models.NewSecretV1Beta(name, envVars, files)

			c := client.NewSecretV1BetaApi()

			_, err = c.CreateSecret(&secret)

			utils.Check(err)

			fmt.Printf("Secret '%s' created.\n", secret.Metadata.Name)
		} else {
			files := parseProjectSecretFiles(fileFlags)
			envVars := parseProjectSecretEnvVars(envFlags)

			secret := models.NewProjectSecretV1(name, envVars, files)

			c := client.NewProjectSecretV1Api(projectID)

			_, err = c.CreateSecret(&secret)

			utils.Check(err)

			fmt.Printf("Project secret '%s' created.\n", secret.Metadata.Name)
		}
	}

	return cmd
}

func parseProjectSecretFiles(fileFlags []string) []models.ProjectSecretV1File {
	var files []models.ProjectSecretV1File

	for _, fileFlag := range fileFlags {
		p, c := parseFile(fileFlag)

		var file models.ProjectSecretV1File
		file.Path = p
		file.Content = encodeFromFileAt(c)
		files = append(files, file)
	}
	return files
}

func parseSecretFiles(fileFlags []string) []models.SecretV1BetaFile {
	var files []models.SecretV1BetaFile

	for _, fileFlag := range fileFlags {
		p, c := parseFile(fileFlag)

		var file models.SecretV1BetaFile
		file.Path = p
		file.Content = encodeFromFileAt(c)
		files = append(files, file)
	}
	return files
}

func parseFile(fileFlag string) (path, content string) {
	matchFormat, err := regexp.MatchString(`^[^: ]+:[^: ]+$`, fileFlag)
	utils.Check(err)

	if !matchFormat {
		utils.Fail("The format of --file flag must be: <local-path>:<semaphore-path>")
	}

	flagPaths := strings.Split(fileFlag, ":")

	return flagPaths[1], flagPaths[0]
}

func parseProjectSecretEnvVars(envFlags []string) []models.ProjectSecretV1EnvVar {
	var envVars []models.ProjectSecretV1EnvVar

	for _, envFlag := range envFlags {
		n, v := parseEnvVar(envFlag)

		var envVar models.ProjectSecretV1EnvVar
		envVar.Name = n
		envVar.Value = v
		envVars = append(envVars, envVar)
	}

	return envVars
}

func parseSecretEnvVars(envFlags []string) []models.SecretV1BetaEnvVar {
	var envVars []models.SecretV1BetaEnvVar

	for _, envFlag := range envFlags {
		n, v := parseEnvVar(envFlag)

		var envVar models.SecretV1BetaEnvVar
		envVar.Name = n
		envVar.Value = v
		envVars = append(envVars, envVar)
	}

	return envVars
}

func parseEnvVar(envFlag string) (name, value string) {
	matchFormat, err := regexp.MatchString(`^.+=.+$`, envFlag)
	utils.Check(err)

	if !matchFormat {
		utils.Fail("The format of -e flag must be: <NAME>=<VALUE>")
	}

	parts := strings.SplitN(envFlag, "=", 2)

	return parts[0], parts[1]
}