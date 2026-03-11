package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPipelinesListV1Alpha_UnmarshalJSON(t *testing.T) {
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
	err := list.UnmarshalJSON([]byte(input))
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

func TestPipelinesListV1Alpha_UnmarshalJSON_Empty(t *testing.T) {
	var list PipelinesListV1Alpha
	err := list.UnmarshalJSON([]byte("[]"))
	assert.Nil(t, err)
	assert.Len(t, list, 0)
}

func TestPipelinesListV1Alpha_UnmarshalJSON_Invalid(t *testing.T) {
	var list PipelinesListV1Alpha
	err := list.UnmarshalJSON([]byte("not json"))
	assert.NotNil(t, err)
}

func TestPipelinesListV1Alpha_MarshalJSON(t *testing.T) {
	list := PipelinesListV1Alpha{
		{
			Id:    "abc-123",
			Name:  "pipeline-1",
			State: "DONE",
			Label: "main",
		},
	}
	list[0].CreatedAt.Seconds = 1700000000

	data, err := list.MarshalJSON()
	assert.Nil(t, err)

	// Round-trip: unmarshal the marshaled output
	var roundTripped PipelinesListV1Alpha
	err = roundTripped.UnmarshalJSON(data)
	assert.Nil(t, err)
	assert.Len(t, roundTripped, 1)
	assert.Equal(t, "abc-123", roundTripped[0].Id)
	assert.Equal(t, "DONE", roundTripped[0].State)
	assert.Equal(t, int64(1700000000), roundTripped[0].CreatedAt.Seconds)
}

func TestPipelinesListV1Alpha_MarshalJSON_Empty(t *testing.T) {
	list := PipelinesListV1Alpha{}
	data, err := list.MarshalJSON()
	assert.Nil(t, err)
	assert.Equal(t, "[]", string(data))
}
