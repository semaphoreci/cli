package deployment_targets

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	client "github.com/semaphoreci/cli/api/client"
	"github.com/semaphoreci/cli/api/models"
	"github.com/semaphoreci/cli/cmd/utils"
	"github.com/spf13/cobra"
)

func Describe(targetId string) {
	c := client.NewDeploymentTargetsV1AlphaApi()

	var deploymentTarget *models.DeploymentTargetV1Alpha
	var err error
	if targetId != "" {
		deploymentTarget, err = c.Describe(targetId)
		utils.Check(err)
	} else {
		utils.Check(errors.New("target id or name must be provided"))
	}
	deploymentTargetYaml, err := deploymentTarget.ToYaml()
	utils.Check(err)

	fmt.Printf("%s\n", deploymentTargetYaml)
}

func DescribeByName(targetName, projectId string) {
	c := client.NewDeploymentTargetsV1AlphaApi()

	deploymentTarget, err := c.DescribeByName(targetName, projectId)
	utils.Check(err)

	deploymentTargetYaml, err := deploymentTarget.ToYaml()
	utils.Check(err)

	fmt.Printf("%s\n", deploymentTargetYaml)
}

func History(targetId string, cmd *cobra.Command) {
	historyRequest := models.HistoryRequestFiltersV1Alpha{
		CursorType: models.HistoryRequestCursorTypeFirstV1Alpha,
	}
	var err error
	historyRequest.GitRefType, err = cmd.Flags().GetString("git-ref-type")
	utils.Check(err)
	historyRequest.GitRefType = strings.ToLower(strings.TrimSpace(historyRequest.GitRefType))

	historyRequest.GitRefLabel, err = cmd.Flags().GetString("git-ref-label")
	utils.Check(err)
	historyRequest.TriggeredBy, err = cmd.Flags().GetString("triggered-by")
	utils.Check(err)

	parameters, err := cmd.Flags().GetStringArray("parameter")
	utils.Check(err)
	for i, parameter := range parameters {
		switch i {
		case 0:
			historyRequest.Parameter1 = parameter
		case 1:
			historyRequest.Parameter2 = parameter
		case 2:
			historyRequest.Parameter3 = parameter
		}
	}
	afterTimestamp, err := cmd.Flags().GetString("after")
	utils.Check(err)
	if afterTimestamp != "" {
		_, err = strconv.ParseInt(afterTimestamp, 10, 64)
		if err != nil {
			utils.Check(errors.New("after timestamp must be valid UNIX time in microseconds"))
		}
		historyRequest.CursorType = models.HistoryRequestCursorTypeAfterV1Alpha
		historyRequest.CursorValue = afterTimestamp
	}

	beforeTimestamp, err := cmd.Flags().GetString("before")
	utils.Check(err)
	if beforeTimestamp != "" {
		_, err = strconv.ParseInt(beforeTimestamp, 10, 64)
		if err != nil {
			utils.Check(errors.New("before timestamp must be valid UNIX time in microseconds"))
		}
		historyRequest.CursorType = models.HistoryRequestCursorTypeBeforeV1Alpha
		historyRequest.CursorValue = beforeTimestamp
	}
	if afterTimestamp != "" && beforeTimestamp != "" {
		utils.Check(errors.New("you can't use both after and before timestamps"))
	}
	c := client.NewDeploymentTargetsV1AlphaApi()

	deployments, err := c.History(targetId, historyRequest)
	utils.Check(err)

	deploymentsYaml, err := deployments.ToYaml()
	utils.Check(err)

	fmt.Printf("%s\n", deploymentsYaml)
}

func HistoryByName(targetName, projectId string, cmd *cobra.Command) {
	c := client.NewDeploymentTargetsV1AlphaApi()

	deploymentTarget, err := c.DescribeByName(targetName, projectId)
	utils.Check(err)

	History(deploymentTarget.Id, cmd)
}

func List(projectId string) {
	c := client.NewDeploymentTargetsV1AlphaApi()

	deploymentTargetsList, err := c.List(projectId)
	utils.Check(err)

	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, "DEPLOYMENT TARGET ID\tNAME\tCREATION TIME (UTC)\tSTATE\tSTATUS")
	if deploymentTargetsList == nil {
		return
	}
	for _, t := range *deploymentTargetsList {
		createdAt := "N/A"
		if t.CreatedAt != nil {
			createdAt = t.CreatedAt.Format("2006-01-02 15:04:05")
		}
		stateName := t.State
		status := "inactive"
		if t.Active {
			status = "active"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", t.Id, t.Name, createdAt, stateName, status)
	}
}
