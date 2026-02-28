package tasks

import (
	"encoding/json"
	"fmt"
	"strings"

	client "github.com/semaphoreci/cli/api/client"
	"github.com/semaphoreci/cli/api/models"
	"github.com/semaphoreci/cli/cmd/utils"
)

func Run(id string, branch string, tag string, pipelineFile string, params []string) {
	req := models.RunTaskRequest{}

	if branch != "" {
		req.Reference = &models.RunTaskReference{
			Type: "BRANCH",
			Name: branch,
		}
	} else if tag != "" {
		req.Reference = &models.RunTaskReference{
			Type: "TAG",
			Name: tag,
		}
	}

	if pipelineFile != "" {
		req.PipelineFile = pipelineFile
	}

	if len(params) > 0 {
		req.Parameters = make(map[string]string)
		for _, p := range params {
			parts := strings.SplitN(p, "=", 2)
			if len(parts) == 2 {
				req.Parameters[parts[0]] = parts[1]
			} else {
				utils.Check(fmt.Errorf("invalid parameter format '%s', expected KEY=VALUE", p))
			}
		}
	}

	body, err := json.Marshal(req)
	utils.Check(err)

	c := client.NewTasksV1AlphaApi()
	resp, err := c.RunTask(id, body)
	utils.Check(err)

	fmt.Printf("Task '%s' triggered. Workflow started: %s\n", id, resp.WorkflowID)
}
