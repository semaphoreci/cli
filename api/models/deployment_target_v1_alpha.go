package models

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"
)

const (
	DeploymentTargetKindV1Alpha = "DeploymentTarget"
)

var DeploymentTargetCmdAliases = []string{"deployment-target", "deployment-targets", "dt", "deployment", "deployments", "deps", "dts"}

type DeploymentTargetListV1Alpha []*DeploymentTargetV1Alpha

type DeploymentTargetV1Alpha struct {
	ApiVersion string `json:"-" yaml:"apiVersion"`
	Kind       string `json:"-" yaml:"kind"`

	DeploymentTargetMetadataV1Alpha `yaml:"metadata,omitempty"`
	DeploymentTargetSpecV1Alpha     `yaml:"spec,omitempty"`
}

type DeploymentTargetMetadataV1Alpha struct {
	Id             string     `json:"id" yaml:"id"`
	Name           string     `json:"name" yaml:"name"`
	ProjectId      string     `json:"project_id" yaml:"project_id"`
	OrganizationId string     `json:"organization_id" yaml:"organization_id"`
	CreatedBy      string     `json:"created_by,omitempty" yaml:"created_by,omitempty"`
	UpdatedBy      string     `json:"updated_by,omitempty" yaml:"updated_by,omitempty"`
	CreatedAt      *time.Time `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
	Description    string     `json:"description" yaml:"description"`
	Url            string     `json:"url" yaml:"url"`
}

type DeploymentTargetSpecV1Alpha struct {
	State              string                          `json:"state" yaml:"state"`
	StateMessage       string                          `json:"state_message" yaml:"state_message"`
	SubjectRules       []*SubjectRuleV1Alpha           `json:"subject_rules" yaml:"subject_rules"`
	ObjectRules        []*ObjectRuleV1Alpha            `json:"object_rules" yaml:"object_rules"`
	LastDeployment     *DeploymentV1Alpha              `json:"last_deployment,omitempty" yaml:"last_deployment,omitempty"`
	Active             bool                            `json:"active" yaml:"active"`
	BookmarkParameter1 string                          `json:"bookmark_parameter1" yaml:"bookmark_parameter1"`
	BookmarkParameter2 string                          `json:"bookmark_parameter2" yaml:"bookmark_parameter2"`
	BookmarkParameter3 string                          `json:"bookmark_parameter3" yaml:"bookmark_parameter3"`
	EnvVars            *DeploymentTargetEnvVarsV1Alpha `json:"env_vars,omitempty" yaml:"env_vars,omitempty"`
	Files              *DeploymentTargetFilesV1Alpha   `json:"files,omitempty" yaml:"files,omitempty"`
}

type DeploymentV1Alpha struct {
	Id             string                    `json:"id" yaml:"id"`
	TargetId       string                    `json:"target_id" yaml:"target_id"`
	PrevPipelineId string                    `json:"prev_pipeline_id" yaml:"prev_pipeline_id"`
	PipelineId     string                    `json:"pipeline_id" yaml:"pipeline_id"`
	TriggeredBy    string                    `json:"triggered_by,omitempty" yaml:"triggered_by,omitempty"`
	TriggeredAt    *time.Time                `json:"triggered_at,omitempty" yaml:"triggered_at,omitempty"`
	State          string                    `json:"state" yaml:"state"`
	StateMessage   string                    `json:"state_message" yaml:"state_message"`
	SwitchId       string                    `json:"switch_id" yaml:"switch_id"`
	TargetName     string                    `json:"target_name" yaml:"target_name"`
	EnvVars        *DeploymentEnvVarsV1Alpha `json:"env_vars,omitempty" yaml:"env_vars,omitempty"`
}
type DeploymentsV1Alpha struct {
	Deployments []DeploymentV1Alpha `json:"deployments,omitempty" yaml:"deployments,omitempty"`
}

type DeploymentEnvVarsV1Alpha []*DeploymentEnvVarV1Alpha

type DeploymentEnvVarV1Alpha struct {
	Name  string `json:"name" yaml:"name"`
	Value string `json:"value" yaml:"value"`
}

type SubjectRuleV1Alpha struct {
	Type      string `json:"type" yaml:"type"`
	SubjectId string `json:"subject_id,omitempty" yaml:"subject_id,omitempty"`
	GitLogin  string `json:"git_login,omitempty" yaml:"git_login,omitempty"`
}

type ObjectRuleV1Alpha struct {
	Type      string `json:"type" yaml:"type"`
	MatchMode string `json:"match_mode" yaml:"match_mode"`
	Pattern   string `json:"pattern" yaml:"pattern"`
}

type HistoryRequest_FiltersV1Alpha struct {
	GitRefType  string `json:"git_ref_type,omitempty" yaml:"git_ref_type,omitempty"`
	GitRefLabel string `json:"git_ref_label,omitempty" yaml:"git_ref_label,omitempty"`
	TriggeredBy string `json:"triggered_by,omitempty" yaml:"triggered_by,omitempty"`
	Parameter1  string `json:"parameter1,omitempty" yaml:"parameter1,omitempty"`
	Parameter2  string `json:"parameter2,omitempty" yaml:"parameter2,omitempty"`
	Parameter3  string `json:"parameter3,omitempty" yaml:"parameter3,omitempty"`
}

type DeploymentTargetEnvVarV1Alpha struct {
	Name  string        `json:"name" yaml:"name"`
	Value HashedContent `json:"value" yaml:"value"`
}

type DeploymentTargetFileV1Alpha struct {
	Path    string        `json:"path" yaml:"path"`
	Content HashedContent `json:"content" yaml:"content"`
	Source  string        `json:"-" yaml:"source"`
}

type HashedContent string

type DeploymentTargetFilesV1Alpha []*DeploymentTargetFileV1Alpha
type DeploymentTargetEnvVarsV1Alpha []*DeploymentTargetEnvVarV1Alpha

type DeploymentTargetCreateRequestV1Alpha struct {
	DeploymentTargetV1Alpha
	UniqueToken string `json:"unique_token" yaml:"-"`
}

type DeploymentTargetUpdateRequestV1Alpha DeploymentTargetCreateRequestV1Alpha

type CordonResponseV1Alpha struct {
	TargetId string `json:"target_id,omitempty" yaml:"target_id,omitempty"`
	Cordoned bool   `json:"cordoned,omitempty" yaml:"cordoned,omitempty"`
}

func NewDeploymentTargetV1AlphaFromJson(data []byte) (*DeploymentTargetV1Alpha, error) {
	dt := DeploymentTargetV1Alpha{}
	err := json.Unmarshal(data, &dt)

	if err != nil {
		return nil, err
	}

	dt.setApiVersionAndKind()

	return &dt, nil
}

func NewDeploymentTargetV1AlphaFromYaml(data []byte) (*DeploymentTargetV1Alpha, error) {
	dt := DeploymentTargetV1Alpha{}

	err := yaml.UnmarshalStrict(data, &dt)

	if err != nil {
		return nil, err
	}

	dt.setApiVersionAndKind()

	return &dt, nil
}

func (dt *DeploymentTargetV1Alpha) ToYaml() ([]byte, error) {
	return yaml.Marshal(dt)
}

func (dt *DeploymentTargetV1Alpha) ToJson() ([]byte, error) {
	return json.Marshal(dt)
}

func (dt *DeploymentTargetV1Alpha) setApiVersionAndKind() {
	dt.ApiVersion = "v1alpha"
	dt.Kind = DeploymentTargetKindV1Alpha
}

func (dt *DeploymentTargetV1Alpha) LoadFiles() error {
	if dt.Files == nil {
		return nil
	}
	return dt.Files.Load()
}

func NewDeploymentTargetListV1AlphaFromJson(data []byte) (*DeploymentTargetListV1Alpha, error) {
	targetList := DeploymentTargetListV1Alpha{}

	err := json.Unmarshal(data, &targetList)

	if err != nil {
		return nil, err
	}
	for i := range targetList {
		targetList[i].setApiVersionAndKind()
	}
	return &targetList, nil
}

func NewCordonResponseV1AlphaFromJson(data []byte) (*CordonResponseV1Alpha, error) {
	cordonResponse := CordonResponseV1Alpha{}

	err := json.Unmarshal(data, &cordonResponse)

	if err != nil {
		return nil, err
	}

	return &cordonResponse, nil
}

func NewDeploymentsV1AlphaFromJson(data []byte) (*DeploymentsV1Alpha, error) {
	deployments := DeploymentsV1Alpha{}

	err := json.Unmarshal(data, &deployments)

	if err != nil {
		return nil, err
	}

	return &deployments, nil
}

func (d *DeploymentsV1Alpha) ToYaml() ([]byte, error) {
	return yaml.Marshal(d)
}

func (c *HashedContent) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	*c = HashedContent(strings.TrimSuffix(string(s), " [md5]"))
	return nil
}

func (c HashedContent) MarshalYAML() (data interface{}, err error) {
	return fmt.Sprintf("%s [md5]", c), nil
}

func (r DeploymentTargetCreateRequestV1Alpha) ObjectName() string {
	return fmt.Sprintf("%s/%s", DeploymentTargetKindV1Alpha, r.Name)
}

func (r DeploymentTargetUpdateRequestV1Alpha) ObjectName() string {
	return fmt.Sprintf("%s/%s", DeploymentTargetKindV1Alpha, r.Name)
}

func (r *DeploymentTargetCreateRequestV1Alpha) ToYaml() ([]byte, error) {
	return yaml.Marshal(r)
}

func (r *DeploymentTargetCreateRequestV1Alpha) ToJson() ([]byte, error) {
	return json.Marshal(r)
}

func (r *DeploymentTargetCreateRequestV1Alpha) LoadFiles() error {
	if r.Files == nil {
		return nil
	}
	return r.Files.Load()
}

func (r *DeploymentTargetUpdateRequestV1Alpha) LoadFiles() error {
	if r.Files == nil {
		return nil
	}
	return r.Files.Load()
}

func (f *DeploymentTargetFilesV1Alpha) Load() error {
	for _, file := range *f {
		if file == nil {
			continue
		}
		if file.Content != "" {
			continue
		}
		err := file.LoadContent()
		if err != nil {
			return err
		}
	}
	return nil
}

func NewDeploymentTargetCreateRequestV1AlphaFromYaml(data []byte) (*DeploymentTargetCreateRequestV1Alpha, error) {
	deploymentTarget := DeploymentTargetCreateRequestV1Alpha{}

	err := yaml.UnmarshalStrict(data, &deploymentTarget)

	if err != nil {
		return nil, err
	}

	deploymentTarget.setApiVersionAndKind()

	return &deploymentTarget, nil
}

func NewDeploymentTargetUpdateRequestV1AlphaFromYaml(data []byte) (*DeploymentTargetUpdateRequestV1Alpha, error) {
	deploymentTargetUpdate := DeploymentTargetUpdateRequestV1Alpha{}

	err := yaml.UnmarshalStrict(data, &deploymentTargetUpdate)

	if err != nil {
		return nil, err
	}

	deploymentTargetUpdate.setApiVersionAndKind()
	if err = deploymentTargetUpdate.LoadFiles(); err != nil {
		return nil, err
	}

	return &deploymentTargetUpdate, nil
}

func (f *DeploymentTargetFileV1Alpha) LoadContent() error {
	content, err := ioutil.ReadFile(f.Source)
	if err != nil {
		return err
	}

	f.Content = HashedContent(base64.StdEncoding.EncodeToString(content))
	return nil
}
