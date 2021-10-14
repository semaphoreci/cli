package workflows

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/semaphoreci/cli/api/client"
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
  defer w.Flush()

  if len(j) == 0 {
	utils.Check(errors.New("failed with empty response, check if the workflow id is correct"))
  }

  fmt.Fprintf(w, "Label: %s\n\n", j[0].Label)
  fmt.Fprintln(w, "PIPELINE ID\tPIPELINE NAME\tCREATION TIME\tSTATE")

  for _, p := range j {
    createdAt := time.Unix(p.CreatedAt.Seconds, 0).Format("2006-01-02 15:04:05")
    fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", p.Id, p.Name, createdAt, p.State)
  }
}
