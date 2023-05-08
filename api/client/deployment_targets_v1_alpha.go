package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	models "github.com/semaphoreci/cli/api/models"
	"github.com/semaphoreci/cli/api/uuid"
)

type DeploymentTargetsApiV1AlphaApi struct {
	BaseClient           BaseClient
	ResourceNameSingular string
	ResourceNamePlural   string
}

func NewDeploymentTargetsV1AlphaApi() DeploymentTargetsApiV1AlphaApi {
	baseClient := NewBaseClientFromConfig()
	baseClient.SetApiVersion("v1alpha")

	return DeploymentTargetsApiV1AlphaApi{
		BaseClient:           baseClient,
		ResourceNamePlural:   "deployment_targets",
		ResourceNameSingular: "deployment_target",
	}
}

func (c *DeploymentTargetsApiV1AlphaApi) Describe(targetId, projectId string) (*models.DeploymentTargetV1Alpha, error) {
	if targetId == "" {
		return nil, errors.New("target id or name must be provided")
	}
	query := targetId
	if projectId != "" {
		query = fmt.Sprintf("%s?project_id=%s", targetId, projectId)
	}
	body, status, err := c.BaseClient.Get(c.ResourceNamePlural, query)
	if err != nil {
		return nil, fmt.Errorf("connecting to Semaphore failed '%s'", err)
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf("http status %d with message \"%s\" received from upstream", status, body)
	}

	return models.NewDeploymentTargetV1AlphaFromJson(body)
}

func (c *DeploymentTargetsApiV1AlphaApi) DescribeByName(targetName, projectId string) (*models.DeploymentTargetV1Alpha, error) {
	if targetName == "" {
		return nil, errors.New("target name must be provided")
	}
	if projectId == "" {
		return nil, errors.New("project id must be provided")
	}
	path := fmt.Sprintf("%s/%s", projectId, targetName)
	body, status, err := c.BaseClient.Get(c.ResourceNamePlural, path)
	if err != nil {
		return nil, fmt.Errorf("connecting to Semaphore failed '%s'", err)
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf("http status %d with message \"%s\" received from upstream", status, body)
	}

	return models.NewDeploymentTargetV1AlphaFromJson(body)
}

func (c *DeploymentTargetsApiV1AlphaApi) History(targetId, projectId string) (*models.DeploymentsV1Alpha, error) {
	if targetId == "" {
		return nil, errors.New("target id must be provided")
	}
	if projectId == "" {
		return nil, errors.New("project id must be provided")
	}
	query := fmt.Sprintf("/%s/history?project_id=%s", targetId, projectId)
	body, status, err := c.BaseClient.Get(c.ResourceNamePlural, query)
	if err != nil {
		return nil, fmt.Errorf("connecting to Semaphore failed '%s'", err)
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf("http status %d with message \"%s\" received from upstream", status, body)
	}

	return models.NewDeploymentsV1AlphaFromJson(body)
}

func (c *DeploymentTargetsApiV1AlphaApi) List(projectId string) (*models.DeploymentTargetListV1Alpha, error) {
	kind := fmt.Sprintf("%s?project_id=%s", c.ResourceNamePlural, projectId)
	body, status, err := c.BaseClient.List(kind)
	if err != nil {
		return nil, fmt.Errorf("connecting to Semaphore failed '%s'", err)
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf("http status %d with message \"%s\" received from upstream", status, body)
	}
	return models.NewDeploymentTargetListV1AlphaFromJson(body)
}

func (c *DeploymentTargetsApiV1AlphaApi) Delete(targetId, projectId string) error {
	unique_token, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("unique token generation failed: %s", err)
	}
	params := fmt.Sprintf("%s?project_id=%s&unique_token=%s", targetId, projectId, unique_token)

	body, status, err := c.BaseClient.Delete(c.ResourceNamePlural, params)
	if err != nil {
		return err
	}

	if status != http.StatusOK {
		return fmt.Errorf("http status %d with message \"%s\" received from upstream", status, body)
	}

	return nil
}

func (c *DeploymentTargetsApiV1AlphaApi) Create(createRequest *models.DeploymentTargetCreateRequestV1Alpha) (*models.DeploymentTargetV1Alpha, error) {
	if createRequest == nil {
		return nil, errors.New("create request must not be nil")
	}
	err := createRequest.LoadFiles()
	if err != nil {
		return nil, err
	}
	unique_token, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("unique token generation failed: %s", err)
	}
	createRequest.UniqueToken = unique_token.String()
	json_body, err := json.Marshal(createRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize deployment target create request: %s", err)
	}

	body, status, err := c.BaseClient.Post(c.ResourceNamePlural, json_body)
	if err != nil {
		return nil, fmt.Errorf("creating %s on Semaphore failed '%s'", c.ResourceNamePlural, err)
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf("http status %d with message \"%s\" received from upstream", status, body)
	}

	return models.NewDeploymentTargetV1AlphaFromJson(body)
}

func (c *DeploymentTargetsApiV1AlphaApi) Update(updateRequest *models.DeploymentTargetUpdateRequestV1Alpha) (*models.DeploymentTargetV1Alpha, error) {
	if updateRequest == nil {
		return nil, errors.New("update request must not be nil")
	}
	if updateRequest.Id == "" {
		return nil, errors.New("update request id must not be empty")
	}
	if updateRequest.ProjectId == "" {
		return nil, errors.New("update request project id must not be empty")
	}
	unique_token, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("unique token generation failed: %s", err)
	}
	updateRequest.UniqueToken = unique_token.String()

	json_body, err := json.Marshal(updateRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize deployment target update request: %s", err)
	}
	body, status, err := c.BaseClient.Patch(c.ResourceNamePlural, updateRequest.Id, json_body)
	if err != nil {
		return nil, fmt.Errorf("creating %s on Semaphore failed '%s'", c.ResourceNamePlural, err)
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf("http status %d with message \"%s\" received from upstream", status, body)
	}

	return models.NewDeploymentTargetV1AlphaFromJson(body)
}

func (c *DeploymentTargetsApiV1AlphaApi) Rebuild(targetId string) (bool, error) {
	return c.cordon(targetId, "on")
}

func (c *DeploymentTargetsApiV1AlphaApi) Stop(targetId string) (bool, error) {
	return c.cordon(targetId, "off")
}

func (c *DeploymentTargetsApiV1AlphaApi) cordon(targetId, opName string) (bool, error) {
	query := fmt.Sprintf("%s/%s", targetId, opName)

	body, status, err := c.BaseClient.Patch(c.ResourceNamePlural, query, nil)
	if err != nil {
		return false, fmt.Errorf("creating %s on Semaphore failed '%s'", c.ResourceNamePlural, err)
	}
	if status != http.StatusOK {
		return false, fmt.Errorf("http status %d with message \"%s\" received from upstream", status, body)
	}

	response, err := models.NewCordonResponseV1AlphaFromJson(body)
	if err != nil {
		return false, fmt.Errorf("wrong response: %s", err)
	}
	if response.TargetId != targetId {
		return false, fmt.Errorf("wrong target id in the response")
	}
	if response.Cordoned != (opName == "off") {
		if response.Cordoned {
			return false, fmt.Errorf("target wasn't rebuilt successfully")
		}
		return false, fmt.Errorf("target wasn't stopped successfully")
	}
	return true, nil
}
