package pipelines

import (
	"fmt"
	"time"

	client "github.com/semaphoreci/cli/api/client"
	"github.com/semaphoreci/cli/cmd/utils"
)

func Describe(id string, follow bool) {
	c := client.NewPipelinesV1AlphaApi()

	for {
		ppl, isDone := describe(c, id)

		fmt.Printf("%s\n", ppl)

		if follow == false || isDone {
			return
		}

		time.Sleep(3 * time.Second)
	}
}

func describe(c client.PipelinesApiV1AlphaApi, id string) ([]byte, bool) {
	pplJ, err := c.DescribePpl(id)
	utils.Check(err)
	pplY, err := pplJ.ToYaml()
	utils.Check(err)

	return pplY, pplJ.IsDone()
}
