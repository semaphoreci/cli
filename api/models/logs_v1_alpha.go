package models

import "encoding/json"

type LogsV1Alpha struct {
	Events []Event `json:"events"`
}

type Event struct {
	Timestamp int32  `json:"timestamp"`
	Type      string `json:"event"`
	Output    string `json:"output"`
	Directive string `json:"directive"`
	ExitCode  int32  `json:"exit_code"`
	JobResult string `json:"job_result"`
}

func NewLogsV1AlphaFromJson(data []byte) (*LogsV1Alpha, error) {
	t := LogsV1Alpha{}
	err := json.Unmarshal(data, &t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}
