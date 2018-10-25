package workflows

import (
	"fmt"
	"log"

	client "github.com/semaphoreci/cli/api/client"
	"github.com/semaphoreci/cli/cmd/utils"
)

func CreateSnapshot(projectName string) {
	projectID := GetProjectId(projectName)
	log.Printf("Project ID: %s\n", projectID)

	c := client.NewWorkflowV1AlphaApi()
	body, err := c.CreateSnapshotWf(projectID)
	utils.Check(err)

	fmt.Printf("PPL_ID: %s\n", string(body))
}
