package workflows

import (
	"fmt"

	"github.com/semaphoreci/cli/api/client"
	"github.com/semaphoreci/cli/cmd/utils"
)

func Rebuild(id string) {
	wfClient := client.NewWorkflowV1AlphaApi()
	body, err := wfClient.Rebuild(id)
	utils.Check(err)
	fmt.Printf("%s\n", string(body))
}
