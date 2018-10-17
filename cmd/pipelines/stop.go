package pipelines

import (
	"fmt"

	client "github.com/semaphoreci/cli/api/client"
	"github.com/semaphoreci/cli/cmd/utils"
)

func Stop(id string) {
	c := client.NewPipelinesV1AlphaApi()
	body, err := c.StopPpl(id)
	utils.Check(err)
	fmt.Printf("%s\n", string(body))
}
