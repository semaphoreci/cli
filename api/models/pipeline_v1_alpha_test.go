package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPipelineFromJsonWithAfterTask(t *testing.T) {
	content := `{
		"pipeline": {
			"ppl_id": "abc-123",
			"name": "Deploy",
			"state": "done",
			"result": "passed",
			"with_after_task": true,
			"after_task_id": "zebra-task-456"
		},
		"blocks": [
			{
				"name": "Build",
				"state": "done",
				"result": "passed",
				"jobs": [{"name": "compile", "job_id": "job-1"}]
			}
		]
	}`

	ppl, err := NewPipelineV1AlphaFromJson([]byte(content))
	assert.Nil(t, err)
	assert.Equal(t, "abc-123", ppl.Pipeline.ID)
	assert.Equal(t, "Deploy", ppl.Pipeline.Name)
	assert.True(t, ppl.Pipeline.WithAfterTask)
	assert.Equal(t, "zebra-task-456", ppl.Pipeline.AfterTaskID)
	assert.Len(t, ppl.Blocks, 1)
}

func TestPipelineFromJsonWithoutAfterTask(t *testing.T) {
	content := `{
		"pipeline": {
			"ppl_id": "abc-123",
			"name": "CI",
			"state": "done",
			"result": "passed"
		},
		"blocks": []
	}`

	ppl, err := NewPipelineV1AlphaFromJson([]byte(content))
	assert.Nil(t, err)
	assert.False(t, ppl.Pipeline.WithAfterTask)
	assert.Empty(t, ppl.Pipeline.AfterTaskID)
}

func TestPipelineToYamlWithAfterTask(t *testing.T) {
	ppl := PipelineV1Alpha{}
	ppl.Pipeline.ID = "abc-123"
	ppl.Pipeline.Name = "Deploy"
	ppl.Pipeline.State = "done"
	ppl.Pipeline.Result = "passed"
	ppl.Pipeline.WithAfterTask = true
	ppl.Pipeline.AfterTaskID = "zebra-task-456"

	yamlBytes, err := ppl.ToYaml()
	assert.Nil(t, err)

	yaml := string(yamlBytes)
	assert.Contains(t, yaml, "with_after_task: true")
	assert.Contains(t, yaml, "after_task_id: zebra-task-456")
}

func TestPipelineToYamlWithoutAfterTask(t *testing.T) {
	ppl := PipelineV1Alpha{}
	ppl.Pipeline.ID = "abc-123"
	ppl.Pipeline.Name = "CI"
	ppl.Pipeline.State = "done"
	ppl.Pipeline.Result = "passed"

	yamlBytes, err := ppl.ToYaml()
	assert.Nil(t, err)

	yaml := string(yamlBytes)
	assert.NotContains(t, yaml, "with_after_task")
	assert.NotContains(t, yaml, "after_task_id")
}
