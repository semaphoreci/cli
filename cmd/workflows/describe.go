package workflows

import (
	"fmt"
	// "time"
	"os"
	"text/tabwriter"
  "encoding/json"

	client "github.com/semaphoreci/cli/api/client"
	"github.com/semaphoreci/cli/api/models"
	"github.com/semaphoreci/cli/cmd/utils"
)

func Describe(projectID, wfID string) {
  c := client.NewPipelinesV1AlphaApi()
  body, err := c.ListPplByWfID(projectID, wfID)
  utils.Check(err)

  prettyPrintPipelineList(body)
}

func prettyPrintPipelineList(jsonList []byte) {
  j := models.PipelinesListV1Alpha{}
	err := json.Unmarshal(jsonList, &j)
  utils.Check(err)

  const padding = 3
  w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)

  fmt.Fprintln(w, "PIPELINE ID\t")

  for _, p := range j {
    fmt.Fprintf(w, "%s\t\n", p.Id)
  }

  w.Flush()
}
