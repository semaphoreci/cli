package tasks

import (
	"fmt"
	"os"
	"text/tabwriter"

	client "github.com/semaphoreci/cli/api/client"
	"github.com/semaphoreci/cli/cmd/utils"
	yaml "gopkg.in/yaml.v2"
)

func Describe(id string) {
	c := client.NewTasksV1AlphaApi()
	task, err := c.DescribeTask(id)
	utils.Check(err)

	y, err := yaml.Marshal(task.Schedule)
	utils.Check(err)

	fmt.Printf("%s", y)

	if len(task.Triggers) > 0 {
		fmt.Println()
		fmt.Println("RECENT TRIGGERS:")

		const padding = 3
		w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)

		fmt.Fprintln(w, "TRIGGERED AT\tSTATUS\tWORKFLOW ID\tBRANCH")

		for _, t := range task.Triggers {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				t.TriggeredAt,
				t.SchedulingStatus,
				t.ScheduledWorkflowID,
				t.Branch,
			)
		}

		err = w.Flush()
		utils.Check(err)
	}
}
