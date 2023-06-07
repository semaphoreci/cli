package cmd

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	client "github.com/semaphoreci/cli/api/client"
	models "github.com/semaphoreci/cli/api/models"
	"github.com/semaphoreci/cli/cmd/utils"
	"github.com/spf13/cobra"
)

var reMatchEnvVarPattern = regexp.MustCompile(`^.+=.+$`)
var reMatchFilePattern = regexp.MustCompile(`^[^: ]+:[^: ]+$`)

func NewCreateDeploymentTargetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "deployment_targets [NAME]",
		Short:   "Create a deployment target.",
		Long:    ``,
		Aliases: models.DeploymentTargetCmdAliases,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			createRequest := models.DeploymentTargetCreateRequestV1Alpha{}
			createRequest.ProjectId = GetProjectID(cmd)

			fileFlags, err := cmd.Flags().GetStringArray("file")
			utils.Check(err)
			for _, fileFlag := range fileFlags {
				if !reMatchFilePattern.MatchString(fileFlag) {
					utils.Fail("The format of --file flag must be: <local-path>:<semaphore-path>")
				}

				flagPaths := strings.Split(fileFlag, ":")

				targetFile := &models.DeploymentTargetFileV1Alpha{
					Source: flagPaths[0],
					Path:   flagPaths[1],
				}
				err := targetFile.LoadContent()
				utils.Check(err)
				if createRequest.Files == nil {
					createRequest.Files = &models.DeploymentTargetFilesV1Alpha{}
				}
				*createRequest.Files = append(*createRequest.Files, targetFile)
			}

			envFlags, err := cmd.Flags().GetStringArray("env")
			utils.Check(err)
			for _, envFlag := range envFlags {
				if !reMatchEnvVarPattern.MatchString(envFlag) {
					utils.Fail("The format of -e flag must be: <NAME>=<VALUE>")
				}

				parts := strings.SplitN(envFlag, "=", 2)
				if createRequest.EnvVars == nil {
					createRequest.EnvVars = &models.DeploymentTargetEnvVarsV1Alpha{}
				}
				*createRequest.EnvVars = append(*createRequest.EnvVars, &models.DeploymentTargetEnvVarV1Alpha{
					Name:  parts[0],
					Value: models.HashedContent(parts[1]),
				})
			}

			c := client.NewDeploymentTargetsV1AlphaApi()

			createRequest.Name = args[0]
			createRequest.Description, err = cmd.Flags().GetString("desc")
			utils.Check(err)
			createRequest.Url, err = cmd.Flags().GetString("url")
			utils.Check(err)
			bookmarks, err := cmd.Flags().GetStringArray("bookmark")
			utils.Check(err)
			for i, bookmark := range bookmarks {
				switch i {
				case 0:
					createRequest.BookmarkParameter1 = bookmark
				case 1:
					createRequest.BookmarkParameter2 = bookmark
				case 2:
					createRequest.BookmarkParameter3 = bookmark
				}
			}

			subjectRulesStrs, err := utils.CSVArrayFlag(cmd, "subject_rules", true)
			utils.Check(err)
			for _, subjectRuleStr := range subjectRulesStrs {
				if len(subjectRuleStr) != 2 {
					utils.Check(fmt.Errorf("invalid subject rule: %s, must be in format TYPE,SUBJECT_ID", subjectRuleStr))
				}
				createRequest.SubjectRules = append(createRequest.SubjectRules, &models.SubjectRuleV1Alpha{
					Type:      subjectRuleStr[0],
					SubjectId: subjectRuleStr[1],
				})
			}

			objectRulesStrs, err := utils.CSVArrayFlag(cmd, "object_rules", true)
			utils.Check(err)
			for _, objectRuleStr := range objectRulesStrs {
				if len(objectRuleStr) != 3 {
					utils.Check(fmt.Errorf("invalid object rule: %s, must be in format TYPE,MODE,PATTERN", objectRuleStr))
				}
				createRequest.ObjectRules = append(createRequest.ObjectRules, &models.ObjectRuleV1Alpha{
					Type:      objectRuleStr[0],
					MatchMode: objectRuleStr[1],
					Pattern:   objectRuleStr[2],
				})
			}

			createdTarget, err := c.Create(&createRequest)
			utils.Check(err)
			if createdTarget == nil {
				utils.Check(errors.New("created target must not be nil"))
				return
			}

			fmt.Printf("Deployment target '%s' ('%s') created.\n", createdTarget.Id, createdTarget.Name)
		},
	}

	cmd.Flags().StringP("desc", "d", "", "Description of deployment target")
	cmd.Flags().StringP("url", "u", "", "URL of deployment target")

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
		"Environment Variables given in the format VAR=VALUE",
	)
	cmd.Flags().StringArrayP("bookmark", "b", []string{}, "Bookmarks for deployment target")
	cmd.Flags().StringArrayP("subject_rules", "s", []string{}, "Subject rules for deployment target")
	cmd.Flags().StringArrayP("object_rules", "o", []string{}, "Object rules for deployment target")
	cmd.Flags().StringP("project-name", "p", "", "project name; if not specified will be inferred from git origin")
	cmd.Flags().StringP("project-id", "i", "", "project id; if not specified will be inferred from git origin")

	return cmd
}
