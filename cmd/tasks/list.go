package tasks

import (
	"fmt"
	"os"
	"text/tabwriter"

	client "github.com/semaphoreci/cli/api/client"
	"github.com/semaphoreci/cli/cmd/utils"
)

func List(projectID string) {
	c := client.NewTasksV1AlphaApi()
	taskList, err := c.ListTasks(projectID)
	utils.Check(err)

	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)

	fmt.Fprintln(w, "ID\tNAME\tSCHEDULED\tBRANCH\tPIPELINE FILE\tSTATUS")

	for _, t := range taskList {
		scheduled := fmt.Sprintf("%t", t.Recurring)

		status := ""
		if t.Paused {
			status = "paused"
		} else if t.Suspended {
			status = "suspended"
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			t.ID,
			t.Name,
			scheduled,
			t.Branch,
			t.PipelineFile,
			status,
		)
	}

	err = w.Flush()
	utils.Check(err)
}
