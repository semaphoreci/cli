package deployment_targets

import (
	"fmt"

	"github.com/semaphoreci/cli/api/client"
	"github.com/semaphoreci/cli/cmd/utils"
)

func Stop(targetId string) {
	client := client.NewDeploymentTargetsV1AlphaApi()
	successful, err := client.Stop(targetId)
	utils.Check(err)
	if successful {
		fmt.Printf("Target [%s] was stopped successfully\n", targetId)
	}
}
