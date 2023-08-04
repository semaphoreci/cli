package models

import (
	"encoding/json"
	"fmt"

	"github.com/semaphoreci/cli/cmd/jobs"
	yaml "gopkg.in/yaml.v3"
)

type JobV1AlphaMetadata struct {
	Name       string      `json:"name,omitempty"`
	Id         string      `json:"id,omitempty"`
	CreateTime json.Number `json:"create_time,omitempty,string" yaml:"create_time,omitempty"`
	UpdateTime json.Number `json:"update_time,omitempty,string" yaml:"update_time,omitempty"`
	StartTime  json.Number `json:"start_time,omitempty,string" yaml:"start_time,omitempty"`
	FinishTime json.Number `json:"finish_time,omitempty,string" yaml:"finish_time,omitempty"`
}

type JobV1AlphaSpecSecret struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
}

type JobV1AlphaSpecFile struct {
	Path    string `json:"path,omitempty" yaml:"path,omitempty"`
	Content string `json:"content,omitempty" yaml:"content,omitempty"`
}

type JobV1AlphaSpecEnvVar struct {
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	Value string `json:"value,omitempty" yaml:"value,omitempty"`
}

type JobV1AlphaAgentImagePullSecret struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
}

type JobV1AlphaAgentMachine struct {
	Type    string `json:"type,omitempty" yaml:"type,omitempty"`
	OsImage string `json:"os_image,omitempty" yaml:"os_image,omitempty"`
}

type JobV1AlphaAgentContainer struct {
	Name    string                 `json:"name,omitempty" yaml:"name,omitempty"`
	Image   string                 `json:"image,omitempty" yaml:"image,omitempty"`
	Command string                 `json:"command,omitempty" yaml:"command,omitempty"`
	EnvVars []JobV1AlphaSpecEnvVar `json:"env_vars,omitempty" yaml:"env_vars,omitempty"`
	Secrets []JobV1AlphaSpecSecret `json:"secrets,omitempty" yaml:"secrets,omitempty"`
}

type JobV1AlphaAgent struct {
	Machine          JobV1AlphaAgentMachine           `json:"machine,omitempty" yaml:"machine,omitempty"`
	Containers       []JobV1AlphaAgentContainer       `json:"containers,omitempty" yaml:"containers,omitempty"`
	ImagePullSecrets []JobV1AlphaAgentImagePullSecret `json:"image_pull_secrets,omitempty" yaml:"image_pull_secrets,omitempty"`
}

type JobV1AlphaSpec struct {
	Agent                  JobV1AlphaAgent        `json:"agent,omitempty" yaml:"agent,omitempty"`
	Files                  []JobV1AlphaSpecFile   `json:"files,omitempty" yaml:"files,omitempty"`
	EnvVars                []JobV1AlphaSpecEnvVar `json:"env_vars,omitempty" yaml:"env_vars,omitempty"`
	Secrets                []JobV1AlphaSpecSecret `json:"secrets,omitempty" yaml:"secrets,omitempty"`
	Commands               []string               `json:"commands,omitempty" yaml:"commands,omitempty"`
	EpilogueAlwaysCommands []string               `json:"epilogue_always_commands,omitempty" yaml:"epilogue_always_commands,omitempty"`
	EpilogueOnPassCommands []string               `json:"epilogue_on_pass_commands,omitempty" yaml:"epilogue_on_pass_commands,omitempty"`
	EpilogueOnFailCommands []string               `json:"epilogue_on_fail_commands,omitempty" yaml:"epilogue_on_fail_commands,omitempty"`
	ProjectId              string                 `json:"project_id,omitempty" yaml:"project_id,omitempty"`
}

type JobV1AlphaStatus struct {
	State  string `json:"state" yaml:"state"`
	Result string `json:"result" yaml:"result"`
	Agent  struct {
		Ip    string `json:"ip" yaml:"ip"`
		Name  string `json:"name" yaml:"name"`
		Ports []struct {
			Name   string `json:"name" yaml:"name"`
			Number int32  `json:"number" yaml:"number"`
		} `json:"ports,omitempty"`
	} `json:"agent,omitempty"`
}

type JobV1Alpha struct {
	ApiVersion string `json:"apiVersion,omitempty" yaml:"apiVersion"`
	Kind       string `json:"kind,omitempty" yaml:"kind"`

	Metadata *JobV1AlphaMetadata `json:"metadata,omitempty"`
	Spec     *JobV1AlphaSpec     `json:"spec,omitempty"`
	Status   *JobV1AlphaStatus   `json:"status,omitempty"`
}

func NewJobV1Alpha(name string) *JobV1Alpha {
	j := JobV1Alpha{}

	j.Metadata = &JobV1AlphaMetadata{}
	j.Metadata.Name = name
	j.setApiVersionAndKind()

	return &j
}

func NewJobV1AlphaFromJson(data []byte) (*JobV1Alpha, error) {
	j := JobV1Alpha{}

	err := json.Unmarshal(data, &j)

	if err != nil {
		return nil, err
	}

	j.setApiVersionAndKind()

	return &j, nil
}

func NewJobV1AlphaFromYaml(data []byte) (*JobV1Alpha, error) {
	j := JobV1Alpha{}

	err := yaml.UnmarshalStrict(data, &j)

	if err != nil {
		return nil, err
	}

	j.setApiVersionAndKind()

	return &j, nil
}

func (j *JobV1Alpha) setApiVersionAndKind() {
	j.ApiVersion = "v1alpha"
	j.Kind = "Job"
}

func (j *JobV1Alpha) ObjectName() string {
	return fmt.Sprintf("Jobs/%s", j.Metadata.Id)
}

func (j *JobV1Alpha) ToJson() ([]byte, error) {
	return json.Marshal(j)
}

func (j *JobV1Alpha) ToYaml() ([]byte, error) {
	return yaml.Marshal(j)
}

func (j *JobV1Alpha) IsSelfHosted() bool {
	return jobs.IsSelfHosted(j.Spec.Agent.Machine.Type)
}

func (j *JobV1Alpha) AgentName() string {
	return j.Status.Agent.Name
}
