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

			parsedSubjectRules, err := utils.CSVArrayFlag(cmd, "subject-rule", true)
			utils.Check(err)

			createSubjectRule := func(parsedSubjRule []string) (*models.SubjectRuleV1Alpha, error) {
				if len(parsedSubjRule) == 0 || len(parsedSubjRule) > 2 {
					return nil, fmt.Errorf("invalid subject rule: %q, must be ANY or AUTO or in format TYPE,SUBJECT", parsedSubjRule)
				}
				rule := &models.SubjectRuleV1Alpha{
					Type: strings.ToUpper(strings.TrimSpace(parsedSubjRule[0])),
				}
				switch rule.Type {
				case "ANY", "AUTO":
					return rule, nil
				case "USER":
					if len(parsedSubjRule) == 2 {
						rule.GitLogin = parsedSubjRule[1]
					} else {
						return nil, errors.New("invalid user subject rule: must be in format USER,GIT_LOGIN")
					}
				case "ROLE":
					if len(parsedSubjRule) == 2 {
						rule.SubjectId = parsedSubjRule[1]
					} else {
						return nil, errors.New(`invalid role subject rule: must be in format ROLE,ROLE_ID`)
					}
				default:
					return nil, fmt.Errorf(`invalid subject rule type: %s, must be one of: ANY, USER, ROLE, AUTO`, rule.Type)
				}
				return rule, nil
			}
			for _, parsedSubjectRule := range parsedSubjectRules {
				subjectRule, err := createSubjectRule(parsedSubjectRule)
				utils.Check(err)

				createRequest.SubjectRules = append(createRequest.SubjectRules, subjectRule)
			}

			objectRulesStrs, err := utils.CSVArrayFlag(cmd, "object-rule", true)
			utils.Check(err)

			createObjectRule := func(parsedObjRule []string) (*models.ObjectRuleV1Alpha, error) {
				if len(parsedObjRule) == 0 || len(parsedObjRule) > 3 {
					return nil, fmt.Errorf("invalid object rule: %q, must be PR or TYPE,ALL or TYPE,MODE,PATTERN", parsedObjRule)
				}
				rule := &models.ObjectRuleV1Alpha{
					Type:      strings.ToUpper(strings.TrimSpace(parsedObjRule[0])),
					MatchMode: models.ObjectRuleMatchModeAllV1Alpha,
				}
				switch rule.Type {
				case models.ObjectRuleTypePullRequestV1Alpha:
					return rule, nil
				case models.ObjectRuleTypeBranchV1Alpha, models.ObjectRuleTypeTagV1Alpha:
					if len(parsedObjRule) == 1 {
						err := fmt.Errorf("invalid object rule: must be %s,ALL or %s,EXACT,<PATTERN> or %s,REGEX,<EXPRESSION>", rule.Type, rule.Type, rule.Type)
						return nil, err
					}
					matchMode := strings.ToUpper(strings.TrimSpace(parsedObjRule[1]))
					switch matchMode {
					case models.ObjectRuleMatchModeAllV1Alpha:
						return rule, nil
					case models.ObjectRuleMatchModeRegexV1Alpha, models.ObjectRuleMatchModeExactV1Alpha:
						if len(parsedObjRule) == 2 {
							if matchMode == models.ObjectRuleMatchModeRegexV1Alpha {
								return nil, fmt.Errorf("invalid object rule: must be %s,REGEX,<EXPRESSION>", rule.Type)
							}
							return nil, fmt.Errorf("invalid object rule: must be %s,EXACT,<PATTERN>", rule.Type)
						}
						rule.MatchMode = matchMode
						rule.Pattern = parsedObjRule[2]
						return rule, nil
					default:
						return nil, fmt.Errorf("invalid object rule match mode: %s, must be ALL, EXACT or REGEX", matchMode)
					}
				default:
					return nil, fmt.Errorf("invalid object rule type: %s, must be BRANCH, TAG or PR", rule.Type)
				}
			}
			for _, parsedObjectRule := range objectRulesStrs {
				objectRule, err := createObjectRule(parsedObjectRule)
				utils.Check(err)

				createRequest.ObjectRules = append(createRequest.ObjectRules, objectRule)
			}

			c := client.NewDeploymentTargetsV1AlphaApi()
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
	cmd.Flags().StringArrayP("subject-rule", "s", []string{}, "Subject rules for deployment target")
	cmd.Flags().StringArrayP("object-rule", "o", []string{}, "Object rules for deployment target")
	cmd.Flags().StringP("project-name", "p", "", "project name; if not specified will be inferred from git origin")
	cmd.Flags().StringP("project-id", "i", "", "project id; if not specified will be inferred from git origin")

	return cmd
}
