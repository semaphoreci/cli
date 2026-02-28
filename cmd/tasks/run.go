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

	// --branch and --tag are mutually exclusive (enforced by cobra);
	// branch takes precedence if both somehow arrive.
	if branch != "" {
		req.Reference = &models.RunTaskReference{
			Type: models.RunTaskRefBranch,
			Name: branch,
		}
	} else if tag != "" {
		req.Reference = &models.RunTaskReference{
			Type: models.RunTaskRefTag,
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

	if resp.WorkflowID == "" {
		utils.Fail("task was triggered but the API returned no workflow ID")
	}

	fmt.Printf("Task '%s' triggered. Workflow started: %s\n", id, resp.WorkflowID)
}
