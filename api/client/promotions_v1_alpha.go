package client

import (
	"encoding/json"
	"errors"
	"fmt"
)

type PromotionsApiV1AlphaApi struct {
	BaseClient           BaseClient
	ResourceNameSingular string
	ResourceNamePlural   string
}

func NewPromotionsV1AlphaApi() PromotionsApiV1AlphaApi {
	baseClient := NewBaseClientFromConfig()
	baseClient.SetApiVersion("v1alpha")

	return PromotionsApiV1AlphaApi{
		BaseClient:           baseClient,
		ResourceNamePlural:   "promotions",
		ResourceNameSingular: "promotion",
	}
}

// Keys of the trigger request that the API reserves for itself; a promotion
// parameter cannot use one of them, every other key is forwarded to the
// promoted pipeline as a parameter.
var reservedPromotionKeys = []string{"pipeline_id", "name", "override", "switch_id", "user_id", "request_token"}

// Trigger starts the promotion `name` of the pipeline `pipelineID`, as the
// "Promote" button does in the UI. `parameters` are the values of a
// parameterized promotion; `override` lets the promotion run on a pipeline that
// failed or is still running.
func (c *PromotionsApiV1AlphaApi) Trigger(pipelineID, name string, override bool, parameters map[string]string) ([]byte, error) {
	request := map[string]interface{}{
		"pipeline_id": pipelineID,
		"name":        name,
		"override":    override,
	}

	for key, value := range parameters {
		for _, reserved := range reservedPromotionKeys {
			if key == reserved {
				return nil, fmt.Errorf("'%s' is reserved by the promotions API and cannot be used as a parameter name", key)
			}
		}

		request[key] = value
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("serializing the promotion request failed '%s'", err)
	}

	body, status, err := c.BaseClient.Post(c.ResourceNamePlural, requestBody)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("connecting to Semaphore failed '%s'", err))
	}

	if status != 200 {
		return nil, errors.New(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	return body, nil
}
