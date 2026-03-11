package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPipelinesListV1Alpha_Unmarshal(t *testing.T) {
	input := `[
		{
			"ppl_id": "abc-123",
			"name": "my-pipeline",
			"state": "DONE",
			"created_at": {"seconds": 1700000000},
			"branch_name": "main"
		},
		{
			"ppl_id": "def-456",
			"name": "my-pipeline",
			"state": "QUEUING",
			"created_at": {"seconds": 1700001000},
			"branch_name": "develop"
		}
	]`

	var list PipelinesListV1Alpha
	err := json.Unmarshal([]byte(input), &list)
	assert.Nil(t, err)
	assert.Len(t, list, 2)

	assert.Equal(t, "abc-123", list[0].Id)
	assert.Equal(t, "my-pipeline", list[0].Name)
	assert.Equal(t, "DONE", list[0].State)
	assert.Equal(t, int64(1700000000), list[0].CreatedAt.Seconds)
	assert.Equal(t, "main", list[0].Label)

	assert.Equal(t, "def-456", list[1].Id)
	assert.Equal(t, "QUEUING", list[1].State)
	assert.Equal(t, "develop", list[1].Label)
}

func TestPipelinesListV1Alpha_Unmarshal_Empty(t *testing.T) {
	var list PipelinesListV1Alpha
	err := json.Unmarshal([]byte("[]"), &list)
	assert.Nil(t, err)
	assert.Len(t, list, 0)
}

func TestPipelinesListV1Alpha_Unmarshal_Invalid(t *testing.T) {
	var list PipelinesListV1Alpha
	err := json.Unmarshal([]byte("not json"), &list)
	assert.NotNil(t, err)
}

func TestPipelinesListV1Alpha_RoundTrip(t *testing.T) {
	list := PipelinesListV1Alpha{
		{
			Id:    "abc-123",
			Name:  "pipeline-1",
			State: "DONE",
			Label: "main",
		},
	}
	list[0].CreatedAt.Seconds = 1700000000

	data, err := json.Marshal(list)
	assert.Nil(t, err)

	var roundTripped PipelinesListV1Alpha
	err = json.Unmarshal(data, &roundTripped)
	assert.Nil(t, err)
	assert.Len(t, roundTripped, 1)
	assert.Equal(t, "abc-123", roundTripped[0].Id)
	assert.Equal(t, "DONE", roundTripped[0].State)
	assert.Equal(t, int64(1700000000), roundTripped[0].CreatedAt.Seconds)
}

func TestPipelinesListV1Alpha_Marshal_Empty(t *testing.T) {
	list := PipelinesListV1Alpha{}
	data, err := json.Marshal(list)
	assert.Nil(t, err)
	assert.Equal(t, "[]", string(data))
}
