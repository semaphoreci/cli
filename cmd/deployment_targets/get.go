package deployment_targets

import (
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	client "github.com/semaphoreci/cli/api/client"
	"github.com/semaphoreci/cli/api/models"
	"github.com/semaphoreci/cli/cmd/utils"
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

func History(targetId string) {
	c := client.NewDeploymentTargetsV1AlphaApi()

	deployments, err := c.History(targetId)
	utils.Check(err)

	deploymentsYaml, err := deployments.ToYaml()
	utils.Check(err)

	fmt.Printf("%s\n", deploymentsYaml)
}

func HistoryByName(targetName, projectId string) {
	c := client.NewDeploymentTargetsV1AlphaApi()

	deploymentTarget, err := c.DescribeByName(targetName, projectId)
	utils.Check(err)

	deployments, err := c.History(deploymentTarget.Id)
	utils.Check(err)

	deploymentsYaml, err := deployments.ToYaml()
	utils.Check(err)

	fmt.Printf("%s\n", deploymentsYaml)
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
