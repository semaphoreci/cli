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

func Describe(targetId, targetName, projectId string) {
	c := client.NewDeploymentTargetsV1AlphaApi()

	var deploymentTarget *models.DeploymentTargetV1Alpha
	var err error
	if targetId != "" {
		deploymentTarget, err = c.Describe(targetId, projectId)
		utils.Check(err)
	} else if targetName != "" {
		deploymentTarget, err = c.DescribeByName(targetName, projectId)
		utils.Check(err)
	} else {
		utils.Check(errors.New("target id or name must be provided"))
	}
	deploymentTargetYaml, err := deploymentTarget.ToYaml()
	utils.Check(err)

	fmt.Printf("%s\n", deploymentTargetYaml)
}

func History(targetId, projectId string) {
	c := client.NewDeploymentTargetsV1AlphaApi()

	deployments, err := c.History(targetId, projectId)
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

	fmt.Fprintln(w, "TARGET ID\tTARGET NAME\tCREATION TIME\tSTATE\tCORDONED")
	if deploymentTargetsList == nil {
		return
	}
	for _, t := range *deploymentTargetsList {
		createdAt := "N/A"
		if t.CreatedAt != nil {
			createdAt = t.CreatedAt.Format("2006-01-02 15:04:05")
		}
		stateName := t.State.Name()
		cordoned := "no"
		if t.Cordoned {
			cordoned = "yes"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%v\n", t.Id, t.Name, createdAt, stateName, cordoned)
	}
}
