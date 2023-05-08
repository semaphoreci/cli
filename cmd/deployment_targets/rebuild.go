package deployment_targets

import (
	"fmt"

	"github.com/semaphoreci/cli/api/client"
	"github.com/semaphoreci/cli/cmd/utils"
)

func Rebuild(targetId string) {
	client := client.NewDeploymentTargetsV1AlphaApi()
	successful, err := client.Rebuild(targetId)
	utils.Check(err)
	if successful {
		fmt.Printf("Target [%s] was rebuilt successfully\n", targetId)
	}
}
