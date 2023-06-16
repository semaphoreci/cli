package client

import (
	"fmt"

	models "github.com/semaphoreci/cli/api/models"
)

type TroubleshootApiV1AlphaApi struct {
	BaseClient           BaseClient
	ResourceNameSingular string
	ResourceNamePlural   string
}

func NewTroubleshootV1AlphaApi() TroubleshootApiV1AlphaApi {
	baseClient := NewBaseClientFromConfig()
	baseClient.SetApiVersion("v1alpha")

	return TroubleshootApiV1AlphaApi{
		BaseClient:           baseClient,
		ResourceNamePlural:   "troubleshoot",
		ResourceNameSingular: "troubleshoot",
	}
}

func (c *TroubleshootApiV1AlphaApi) TroubleshootWorkflow(workflowID string) (*models.TroubleshootV1Alpha, error) {
	urlEncode := fmt.Sprintf("%s/workflow", c.ResourceNamePlural)
	body, status, err := c.BaseClient.Get(urlEncode, workflowID)

	if err != nil {
		return nil, fmt.Errorf("connecting to Semaphore failed '%s'", err)
	}

	if status != 200 {
		return nil, fmt.Errorf("http status %d with message \"%s\" received from upstream", status, body)
	}

	return models.NewTroubleshootV1AlphaFromJson(body)
}

func (c *TroubleshootApiV1AlphaApi) TroubleshootJob(jobID string) (*models.TroubleshootV1Alpha, error) {
	urlEncode := fmt.Sprintf("%s/job", c.ResourceNamePlural)
	body, status, err := c.BaseClient.Get(urlEncode, jobID)

	if err != nil {
		return nil, fmt.Errorf("connecting to Semaphore failed '%s'", err)
	}

	if status != 200 {
		return nil, fmt.Errorf("http status %d with message \"%s\" received from upstream", status, body)
	}

	return models.NewTroubleshootV1AlphaFromJson(body)
}

func (c *TroubleshootApiV1AlphaApi) TroubleshootPipeline(pplID string) (*models.TroubleshootV1Alpha, error) {
	urlEncode := fmt.Sprintf("%s/pipeline/", c.ResourceNamePlural)
	body, status, err := c.BaseClient.Get(urlEncode, pplID)

	if err != nil {
		return nil, fmt.Errorf("connecting to Semaphore failed '%s'", err)
	}

	if status != 200 {
		return nil, fmt.Errorf("http status %d with message \"%s\" received from upstream", status, body)
	}

	return models.NewTroubleshootV1AlphaFromJson(body)
}
