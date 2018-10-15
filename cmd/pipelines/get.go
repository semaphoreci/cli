package pipelines

import (
	"fmt"

	client "github.com/semaphoreci/cli/api/client"
	"github.com/semaphoreci/cli/cmd/utils"
)

func Describe(id string) {
	c := client.NewPipelinesV1AlphaApi()
	ppl, err := c.DescribePpl(id)
	utils.Check(err)
	pplY, err := ppl.ToYaml()
	utils.Check(err)
	fmt.Printf("%s\n", pplY)
}
