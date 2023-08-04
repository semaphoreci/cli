package pipelines

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	client "github.com/semaphoreci/cli/api/client"
	"github.com/semaphoreci/cli/api/models"
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

func List(projectID string) {
	fmt.Printf("%s\n", projectID)
	c := client.NewPipelinesV1AlphaApi()
	body, err := c.ListPpl(projectID)
	utils.Check(err)

	prettyPrintPipelineList(body)
}

func prettyPrintPipelineList(jsonList []byte) {
	j := models.PipelinesListV1Alpha{}
	err := json.Unmarshal(jsonList, &j)
	utils.Check(err)

	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)

	fmt.Fprintln(w, "PIPELINE ID\tPIPELINE NAME\tCREATION TIME\tSTATE\tLABEL")

	for _, p := range j {
		createdAt := time.Unix(p.CreatedAt.Seconds, 0).Format("2006-01-02 15:04:05")
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", p.Id, p.Name, createdAt, p.State, p.Label)
	}

	if err := w.Flush(); err != nil {
		fmt.Printf("Error flushing when pretty printing pipelines: %v\n", err)
	}
}
