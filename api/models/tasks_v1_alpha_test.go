package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTaskListV1AlphaFromJSON(t *testing.T) {
	input := `[
		{
			"id": "bb2ba294-d4b3-48bc-90a7-12dd56e9424c",
			"name": "deploy",
			"project_id": "aa1ba294-d4b3-48bc-90a7-12dd56e9424a",
			"branch": "main",
			"at": "0 0 * * *",
			"pipeline_file": ".semaphore/deploy.yml",
			"recurring": true,
			"paused": false,
			"suspended": false
		},
		{
			"id": "cc3ba294-d4b3-48bc-90a7-12dd56e9424d",
			"name": "nightly",
			"project_id": "aa1ba294-d4b3-48bc-90a7-12dd56e9424a",
			"branch": "main",
			"at": "0 2 * * *",
			"pipeline_file": ".semaphore/nightly.yml",
			"recurring": true,
			"paused": true
		}
	]`

	list, err := NewTaskListV1AlphaFromJSON([]byte(input))
	assert.Nil(t, err)
	assert.Len(t, list, 2)

	assert.Equal(t, "bb2ba294-d4b3-48bc-90a7-12dd56e9424c", list[0].ID)
	assert.Equal(t, "deploy", list[0].Name)
	assert.Equal(t, "aa1ba294-d4b3-48bc-90a7-12dd56e9424a", list[0].ProjectID)
	assert.Equal(t, "main", list[0].Branch)
	assert.Equal(t, "0 0 * * *", list[0].At)
	assert.Equal(t, ".semaphore/deploy.yml", list[0].PipelineFile)
	assert.True(t, list[0].Recurring)
	assert.False(t, list[0].Paused)
	assert.False(t, list[0].Suspended)

	assert.Equal(t, "cc3ba294-d4b3-48bc-90a7-12dd56e9424d", list[1].ID)
	assert.Equal(t, "nightly", list[1].Name)
	assert.True(t, list[1].Paused)
}

func TestNewTaskListV1AlphaFromJSON__EmptyList(t *testing.T) {
	list, err := NewTaskListV1AlphaFromJSON([]byte("[]"))
	assert.Nil(t, err)
	assert.Len(t, list, 0)
}

func TestNewTaskListV1AlphaFromJSON__InvalidJSON(t *testing.T) {
	_, err := NewTaskListV1AlphaFromJSON([]byte("not json"))
	assert.NotNil(t, err)
}

func TestNewTaskDescribeV1AlphaFromJSON(t *testing.T) {
	input := `{
		"schedule": {
			"id": "bb2ba294-d4b3-48bc-90a7-12dd56e9424c",
			"name": "deploy",
			"project_id": "aa1ba294-d4b3-48bc-90a7-12dd56e9424a",
			"branch": "main",
			"at": "0 0 * * *",
			"pipeline_file": ".semaphore/deploy.yml",
			"recurring": true,
			"description": "Daily deploy task",
			"parameters": {
				"ENV": "production"
			}
		},
		"triggers": [
			{
				"triggered_at": "2024-01-15 09:00:00",
				"scheduling_status": "passed",
				"scheduled_workflow_id": "dd4ba294-d4b3-48bc-90a7-12dd56e9424e",
				"branch": "main",
				"pipeline_file": ".semaphore/deploy.yml"
			},
			{
				"triggered_at": "2024-01-14 09:00:00",
				"scheduling_status": "failed",
				"scheduled_workflow_id": "ee5ba294-d4b3-48bc-90a7-12dd56e9424f",
				"branch": "main",
				"error_description": "pipeline file not found"
			}
		]
	}`

	desc, err := NewTaskDescribeV1AlphaFromJSON([]byte(input))
	assert.Nil(t, err)

	assert.Equal(t, "bb2ba294-d4b3-48bc-90a7-12dd56e9424c", desc.Schedule.ID)
	assert.Equal(t, "deploy", desc.Schedule.Name)
	assert.Equal(t, "Daily deploy task", desc.Schedule.Description)
	assert.Equal(t, "production", desc.Schedule.Parameters["ENV"])
	assert.True(t, desc.Schedule.Recurring)

	assert.Len(t, desc.Triggers, 2)
	assert.Equal(t, "2024-01-15 09:00:00", desc.Triggers[0].TriggeredAt)
	assert.Equal(t, "passed", desc.Triggers[0].SchedulingStatus)
	assert.Equal(t, "dd4ba294-d4b3-48bc-90a7-12dd56e9424e", desc.Triggers[0].ScheduledWorkflowID)
	assert.Equal(t, "main", desc.Triggers[0].Branch)

	assert.Equal(t, "failed", desc.Triggers[1].SchedulingStatus)
	assert.Equal(t, "pipeline file not found", desc.Triggers[1].ErrorDescription)
}

func TestNewTaskDescribeV1AlphaFromJSON__NoTriggers(t *testing.T) {
	input := `{
		"schedule": {
			"id": "bb2ba294-d4b3-48bc-90a7-12dd56e9424c",
			"name": "deploy",
			"project_id": "aa1ba294-d4b3-48bc-90a7-12dd56e9424a",
			"pipeline_file": ".semaphore/deploy.yml"
		}
	}`

	desc, err := NewTaskDescribeV1AlphaFromJSON([]byte(input))
	assert.Nil(t, err)
	assert.Equal(t, "deploy", desc.Schedule.Name)
	assert.Nil(t, desc.Triggers)
}

func TestNewTaskDescribeV1AlphaFromJSON__InvalidJSON(t *testing.T) {
	_, err := NewTaskDescribeV1AlphaFromJSON([]byte("{invalid"))
	assert.NotNil(t, err)
}

func TestNewRunTaskResponseFromJSON(t *testing.T) {
	input := `{"workflow_id": "dd4ba294-d4b3-48bc-90a7-12dd56e9424e"}`

	resp, err := NewRunTaskResponseFromJSON([]byte(input))
	assert.Nil(t, err)
	assert.Equal(t, "dd4ba294-d4b3-48bc-90a7-12dd56e9424e", resp.WorkflowID)
}

func TestNewRunTaskResponseFromJSON__InvalidJSON(t *testing.T) {
	_, err := NewRunTaskResponseFromJSON([]byte("bad"))
	assert.NotNil(t, err)
}

func TestRunTaskRequest__MarshalBranchReference(t *testing.T) {
	req := RunTaskRequest{
		Reference: &RunTaskReference{
			Type: "BRANCH",
			Name: "main",
		},
		PipelineFile: ".semaphore/custom.yml",
		Parameters: map[string]string{
			"ENV":    "staging",
			"REGION": "us-east-1",
		},
	}

	data, err := json.Marshal(req)
	assert.Nil(t, err)

	var parsed map[string]interface{}
	err = json.Unmarshal(data, &parsed)
	assert.Nil(t, err)

	ref := parsed["reference"].(map[string]interface{})
	assert.Equal(t, "BRANCH", ref["type"])
	assert.Equal(t, "main", ref["name"])
	assert.Equal(t, ".semaphore/custom.yml", parsed["pipeline_file"])

	params := parsed["parameters"].(map[string]interface{})
	assert.Equal(t, "staging", params["ENV"])
	assert.Equal(t, "us-east-1", params["REGION"])
}

func TestRunTaskRequest__MarshalTagReference(t *testing.T) {
	req := RunTaskRequest{
		Reference: &RunTaskReference{
			Type: "TAG",
			Name: "v1.0",
		},
	}

	data, err := json.Marshal(req)
	assert.Nil(t, err)

	var parsed map[string]interface{}
	err = json.Unmarshal(data, &parsed)
	assert.Nil(t, err)

	ref := parsed["reference"].(map[string]interface{})
	assert.Equal(t, "TAG", ref["type"])
	assert.Equal(t, "v1.0", ref["name"])
}

func TestRunTaskRequest__MarshalEmpty(t *testing.T) {
	req := RunTaskRequest{}

	data, err := json.Marshal(req)
	assert.Nil(t, err)
	assert.Equal(t, "{}", string(data))
}

func TestTaskV1Alpha__OmitsEmptyFields(t *testing.T) {
	task := TaskV1Alpha{
		ID:           "abc-123",
		Name:         "deploy",
		ProjectID:    "prj-456",
		PipelineFile: ".semaphore/deploy.yml",
	}

	data, err := json.Marshal(task)
	assert.Nil(t, err)

	var parsed map[string]interface{}
	err = json.Unmarshal(data, &parsed)
	assert.Nil(t, err)

	assert.Equal(t, "abc-123", parsed["id"])
	assert.Equal(t, "deploy", parsed["name"])
	assert.Equal(t, "prj-456", parsed["project_id"])
	assert.Equal(t, ".semaphore/deploy.yml", parsed["pipeline_file"])

	_, hasBranch := parsed["branch"]
	assert.False(t, hasBranch)
	_, hasAt := parsed["at"]
	assert.False(t, hasAt)
	_, hasDescription := parsed["description"]
	assert.False(t, hasDescription)
	_, hasParams := parsed["parameters"]
	assert.False(t, hasParams)
}
