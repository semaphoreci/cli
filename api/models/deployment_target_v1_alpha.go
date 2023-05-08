package models

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"
)

type DeploymentTargetState int32

type SubjectRuleType int32

// Describes the state of deployment.
type DeploymentState int32

// Target object rule type - indicates what kind of workflows are allowed to use deployment target.
type ObjectRuleV1Alpha_Type int32

// Target object rule mode - indicates how pattern is matched against provided workflow metadata.
type ObjectRuleV1Alpha_Mode int32

const (
	DeploymentTargetKindV1Alpha = "DeploymentTarget"

	DeploymentTarget_SYNCING  DeploymentTargetState = 0
	DeploymentTarget_USABLE   DeploymentTargetState = 1
	DeploymentTarget_UNUSABLE DeploymentTargetState = 2
	DeploymentTarget_CORDONED DeploymentTargetState = 3

	SubjectRule_USER  SubjectRuleType = 0
	SubjectRule_ROLE  SubjectRuleType = 1
	SubjectRule_GROUP SubjectRuleType = 2
	SubjectRule_AUTO  SubjectRuleType = 3
	SubjectRule_ANY   SubjectRuleType = 4

	// deployment is being processed
	Deployment_PENDING DeploymentState = 0
	// deployment was processed and pipeline was scheduled
	Deployment_STARTED DeploymentState = 1
	// deployment was processed and pipeline was not scheduled
	Deployment_FAILED DeploymentState = 2

	// workflows triggered by push to a branch
	ObjectRuleV1Alpha_BRANCH ObjectRuleV1Alpha_Type = 0
	// workflows triggered by a new tag
	ObjectRuleV1Alpha_TAG ObjectRuleV1Alpha_Type = 1
	// workflows triggered by a pull request
	ObjectRuleV1Alpha_PR ObjectRuleV1Alpha_Type = 2

	// pattern is discarded, matches all workflows of type
	ObjectRuleV1Alpha_ALL ObjectRuleV1Alpha_Mode = 0
	// pattern is matched exactly with workflow metadata
	ObjectRuleV1Alpha_EXACT ObjectRuleV1Alpha_Mode = 1
	// pattern is compiled to regular expression
	ObjectRuleV1Alpha_REGEX ObjectRuleV1Alpha_Mode = 2
)

// Enum value maps for ObjectRuleV1Alpha_Mode.
var (
	ObjectRuleV1Alpha_Mode_name = map[int32]string{
		0: "ALL",
		1: "EXACT",
		2: "REGEX",
	}
	ObjectRuleV1Alpha_Mode_value = map[string]int32{
		"ALL":   0,
		"EXACT": 1,
		"REGEX": 2,
	}
	ObjectRuleV1Alpha_Type_name = map[int32]string{
		0: "BRANCH",
		1: "TAG",
		2: "PR",
	}
	ObjectRuleV1Alpha_Type_value = map[string]int32{
		"BRANCH": 0,
		"TAG":    1,
		"PR":     2,
	}
	DeploymentState_name = map[int32]string{
		0: "PENDING",
		1: "STARTED",
		2: "FAILED",
	}
	DeploymentState_value = map[string]int32{
		"PENDING": 0,
		"STARTED": 1,
		"FAILED":  2,
	}
	SubjectRuleType_name = map[int32]string{
		0: "USER",
		1: "ROLE",
		2: "GROUP",
		3: "AUTO",
		4: "ANY",
	}
	SubjectRuleType_value = map[string]int32{
		"USER":  0,
		"ROLE":  1,
		"GROUP": 2,
		"AUTO":  3,
		"ANY":   4,
	}
	DeploymentTargetState_name = map[int32]string{
		0: "SYNCING",
		1: "USABLE",
		2: "UNUSABLE",
		3: "CORDONED",
	}
	DeploymentTargetState_value = map[string]int32{
		"SYNCING":  0,
		"USABLE":   1,
		"UNUSABLE": 2,
		"CORDONED": 3,
	}
)

type DeploymentTargetListV1Alpha []*DeploymentTargetV1Alpha

type DescribeResponseV1Alpha struct {
	Target *DeploymentTargetV1Alpha `json:"target,omitempty" yaml:"target,omitempty"`
}

type HistoryResponseV1Alpha struct {
	Deployments  []*DeploymentV1Alpha `json:"deployments,omitempty" yaml:"deployments,omitempty"`
	CursorBefore uint64               `json:"cursor_before,omitempty" yaml:"cursor_before,omitempty"`
	CursorAfter  uint64               `json:"cursor_after,omitempty" yaml:"cursor_after,omitempty"`
}

type CordonResponseV1Alpha struct {
	TargetId string `json:"target_id,omitempty" yaml:"target_id,omitempty"`
	Cordoned bool   `json:"cordoned,omitempty" yaml:"cordoned,omitempty"`
}

type CreateResponseV1Alpha struct {
	Target *DeploymentTargetV1Alpha `json:"target,omitempty" yaml:"target,omitempty"`
}

type UpdateResponseV1Alpha struct {
	Target *DeploymentTargetV1Alpha `json:"target,omitempty" yaml:"target,omitempty"`
}

type DeleteResponseV1Alpha struct {
	TargetId string `json:"target_id,omitempty" yaml:"target_id,omitempty"`
}

type DeploymentTargetV1Alpha struct {
	ApiVersion string `json:"apiVersion,omitempty" yaml:"apiVersion"`
	Kind       string `json:"kind,omitempty" yaml:"kind"`

	DeploymentTargetMetadataV1Alpha `yaml:"metadata,omitempty"`
	DeploymentTargetSpecV1Alpha     `yaml:"spec,omitempty"`
}

type DeploymentTargetMetadataV1Alpha struct {
	Id             string     `json:"id,omitempty" yaml:"id,omitempty"`
	Name           string     `json:"name,omitempty" yaml:"name,omitempty"`
	ProjectId      string     `json:"project_id,omitempty" yaml:"project_id,omitempty"`
	OrganizationId string     `json:"organization_id,omitempty" yaml:"organization_id,omitempty"`
	CreatedBy      string     `json:"created_by,omitempty" yaml:"created_by,omitempty"`
	UpdatedBy      string     `json:"updated_by,omitempty" yaml:"updated_by,omitempty"`
	CreatedAt      *time.Time `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
	Description    string     `json:"description,omitempty" yaml:"description,omitempty"`
	Url            string     `json:"url,omitempty" yaml:"url,omitempty"`
}

type DeploymentTargetSpecV1Alpha struct {
	State              DeploymentTargetState `json:"state" yaml:"state"`
	StateMessage       string                `json:"state_message,omitempty" yaml:"state_message,omitempty"`
	SubjectRules       []*SubjectRuleV1Alpha `json:"subject_rules,omitempty" yaml:"subject_rules,omitempty"`
	ObjectRules        []*ObjectRuleV1Alpha  `json:"object_rules,omitempty" yaml:"object_rules,omitempty"`
	LastDeployment     *DeploymentV1Alpha    `json:"last_deployment,omitempty" yaml:"last_deployment,omitempty"`
	Cordoned           bool                  `json:"cordoned,omitempty" yaml:"cordoned,omitempty"`
	BookmarkParameter1 string                `json:"bookmark_parameter1,omitempty" yaml:"bookmark_parameter1,omitempty"`
	BookmarkParameter2 string                `json:"bookmark_parameter2,omitempty" yaml:"bookmark_parameter2,omitempty"`
	BookmarkParameter3 string                `json:"bookmark_parameter3,omitempty" yaml:"bookmark_parameter3,omitempty"`
	SecretName         string                `json:"secret_name,omitempty" yaml:"secret_name,omitempty"`
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

func NewDeploymentTargetListV1AlphaFromJson(data []byte) (*DeploymentTargetListV1Alpha, error) {
	j := DeploymentTargetListV1Alpha{}

	err := json.Unmarshal(data, &j)

	if err != nil {
		return nil, err
	}

	return &j, nil
}

func NewCordonResponseV1AlphaFromJson(data []byte) (*CordonResponseV1Alpha, error) {
	j := CordonResponseV1Alpha{}

	err := json.Unmarshal(data, &j)

	if err != nil {
		return nil, err
	}

	return &j, nil
}

type DeploymentV1Alpha struct {
	Id             string                           `json:"id,omitempty" yaml:"id,omitempty"`
	TargetId       string                           `json:"target_id,omitempty" yaml:"target_id,omitempty"`
	PrevPipelineId string                           `json:"prev_pipeline_id,omitempty" yaml:"prev_pipeline_id,omitempty"`
	PipelineId     string                           `json:"pipeline_id,omitempty" yaml:"pipeline_id,omitempty"`
	TriggeredBy    string                           `json:"triggered_by,omitempty" yaml:"triggered_by,omitempty"`
	TriggeredAt    *time.Time                       `json:"triggered_at,omitempty" yaml:"triggered_at,omitempty"`
	State          DeploymentState                  `json:"state,omitempty" yaml:"state,omitempty"`
	StateMessage   string                           `json:"state_message,omitempty" yaml:"state_message,omitempty"`
	SwitchId       string                           `json:"switch_id,omitempty" yaml:"switch_id,omitempty"`
	TargetName     string                           `json:"target_name,omitempty" yaml:"target_name,omitempty"`
	EnvVars        []*DeploymentTargetEnvVarV1Alpha `json:"env_vars,omitempty" yaml:"env_vars,omitempty"`
}

type DeploymentsV1Alpha struct {
	Deployments []DeploymentV1Alpha `json:"deployments,omitempty" yaml:"deployments,omitempty"`
}

func NewDeploymentsV1AlphaFromJson(data []byte) (*DeploymentsV1Alpha, error) {
	d := DeploymentsV1Alpha{}

	err := json.Unmarshal(data, &d)

	if err != nil {
		return nil, err
	}

	return &d, nil
}

func (d *DeploymentsV1Alpha) ToYaml() ([]byte, error) {
	return yaml.Marshal(d)
}

type SubjectRuleV1Alpha struct {
	Type      SubjectRuleType `json:"type" yaml:"type"`
	SubjectId string          `json:"subject_id" yaml:"subject_id"`
}

type ObjectRuleV1Alpha struct {
	Type      ObjectRuleV1Alpha_Type `json:"type" yaml:"type"`
	MatchMode ObjectRuleV1Alpha_Mode `json:"match_mode" yaml:"match_mode"`
	Pattern   string                 `json:"pattern" yaml:"pattern"`
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
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	Value string `json:"value,omitempty" yaml:"value,omitempty"`
}

type DeploymentTargetFileV1Alpha struct {
	Path    string `json:"path" yaml:"path"`
	Content string `json:"content" yaml:"content"`
	Source  string `json:"-" yaml:"source"`
}

type DeploymentTargetFiles []*DeploymentTargetFileV1Alpha

type DeploymentTargetCreateRequestV1Alpha struct {
	DeploymentTargetV1Alpha `json:"target" yaml:",inline"`
	UniqueToken             string                           `json:"unique_token" yaml:"-"`
	EnvVars                 []*DeploymentTargetEnvVarV1Alpha `json:"env_vars" yaml:"env_vars"`
	Files                   DeploymentTargetFiles            `json:"files" yaml:"files"`
	ProjectId               string                           `json:"project_id" yaml:"-"`
}

type DeploymentTargetUpdateRequestV1Alpha DeploymentTargetCreateRequestV1Alpha

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
	return r.Files.Load()
}

func (r *DeploymentTargetUpdateRequestV1Alpha) LoadFiles() error {
	return r.Files.Load()
}

func (f *DeploymentTargetFiles) Load() error {
	for _, file := range *f {
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
	r := DeploymentTargetCreateRequestV1Alpha{}

	err := yaml.UnmarshalStrict(data, &r)

	if err != nil {
		return nil, err
	}

	r.setApiVersionAndKind()

	return &r, nil
}

func NewDeploymentTargetUpdateRequestV1AlphaFromYaml(data []byte) (*DeploymentTargetUpdateRequestV1Alpha, error) {
	r := DeploymentTargetUpdateRequestV1Alpha{}

	err := yaml.UnmarshalStrict(data, &r)

	if err != nil {
		return nil, err
	}

	r.setApiVersionAndKind()
	if err = r.LoadFiles(); err != nil {
		return nil, err
	}

	return &r, nil
}

func (f *DeploymentTargetFileV1Alpha) LoadContent() error {
	content, err := ioutil.ReadFile(f.Source)
	if err != nil {
		return err
	}

	f.Content = base64.StdEncoding.EncodeToString(content)
	return nil
}

type unmarshaller func(data []byte, v interface{}) error

func (s *DeploymentTargetState) unmarshal(fun unmarshaller, data []byte) error {
	var val int32
	err := fun(data, &val)
	if err == nil {
		*s = DeploymentTargetState(val)
		return nil
	}
	var str string
	err = fun(data, &str)
	if err != nil {
		return err
	}
	str = strings.TrimSpace(strings.ToUpper(str))

	if val32, ok := DeploymentTargetState_value[str]; !ok {
		return errors.New("invalid state value")
	} else {
		*s = DeploymentTargetState(val32)
	}
	return nil
}

func (s *DeploymentTargetState) UnmarshalJSON(data []byte) error {
	return s.unmarshal(json.Unmarshal, data)
}

func (s *DeploymentTargetState) UnmarshalYAML(data []byte) error {
	return s.unmarshal(yaml.Unmarshal, data)
}

func (s *DeploymentTargetState) MarshalYAML() (interface{}, error) {
	return s.Name(), nil
}

func (s *DeploymentTargetState) MarshalJSON() ([]byte, error) {
	return json.Marshal(*s)
}

func (s *DeploymentTargetState) Name() string {
	if name, ok := DeploymentTargetState_name[int32(*s)]; ok {
		return name
	}
	return "N/A"
}

func ParseSubjectRuleType(valI interface{}) (subjectRuleType SubjectRuleType, err error) {
	val, err := parseEnumValue("subject rule type", valI, SubjectRuleType_value, SubjectRuleType_name)
	return SubjectRuleType(val), err
}

func ParseDeploymentTargetState(valI interface{}) (deploymentTargetState DeploymentState, err error) {
	val, err := parseEnumValue("deployment target state", valI, DeploymentTargetState_value, DeploymentTargetState_name)
	return DeploymentState(val), err
}

func ParseObjectRuleType(valI interface{}) (objectRuleType ObjectRuleV1Alpha_Type, err error) {
	val, err := parseEnumValue("object rule type", valI, ObjectRuleV1Alpha_Type_value, ObjectRuleV1Alpha_Type_name)
	return ObjectRuleV1Alpha_Type(val), err
}

func ParseObjectRuleMode(valI interface{}) (objectRuleMode ObjectRuleV1Alpha_Mode, err error) {
	val, err := parseEnumValue("object rule mode", valI, ObjectRuleV1Alpha_Mode_value, ObjectRuleV1Alpha_Mode_name)
	return ObjectRuleV1Alpha_Mode(val), err
}

func parseEnumValue(enumName string, valI interface{}, nameToValue map[string]int32,
	valueToName map[int32]string) (value int32, err error) {
	value = -1
	switch v := valI.(type) {
	case string:
		i, errParse := strconv.ParseInt(v, 10, 32)
		if errParse != nil {
			val32, ok := SubjectRuleType_value[strings.ToUpper(strings.TrimSpace(v))]
			if !ok {
				return -1, fmt.Errorf("%s value '%v' is invalid", enumName, valI)
			}
			return val32, nil
		}
		value = int32(i)
	case int:
		value = int32(v)
	case int32:
		value = v
	case int64:
		value = int32(v)
	}
	_, ok := valueToName[value]
	if !ok {
		return -1, fmt.Errorf("%s value '%v' is invalid", enumName, valI)
	}
	return
}

func (s *SubjectRuleType) unmarshal(data []byte) error {
	str := strings.ToUpper(strings.TrimSpace(string(data)))
	if str == "" {
		*s = 0
		return nil
	}
	val, err := strconv.ParseInt(str, 10, 32)
	if err == nil {
		*s = SubjectRuleType(val)
		return nil
	}
	val32, ok := SubjectRuleType_value[str]
	if !ok {
		return errors.New("invalid state value")
	}
	*s = SubjectRuleType(val32)
	return nil
}

func (s *SubjectRuleType) UnmarshalJSON(data []byte) error {
	return s.unmarshal(data)
}

func (s *SubjectRuleType) UnmarshalYAML(data []byte) error {
	return s.unmarshal(data)
}

func (s *ObjectRuleV1Alpha_Type) unmarshal(data []byte) error {
	str := strings.ToUpper(strings.TrimSpace(string(data)))
	if str == "" {
		*s = 0
		return nil
	}
	val, err := strconv.ParseInt(str, 10, 32)
	if err == nil {
		*s = ObjectRuleV1Alpha_Type(val)
		return nil
	}
	val32, ok := ObjectRuleV1Alpha_Type_value[str]
	if !ok {
		return errors.New("invalid state value")
	}
	*s = ObjectRuleV1Alpha_Type(val32)
	return nil
}

func (s *ObjectRuleV1Alpha_Type) UnmarshalJSON(data []byte) error {
	return s.unmarshal(data)
}

func (s *ObjectRuleV1Alpha_Type) UnmarshalYAML(data []byte) error {
	return s.unmarshal(data)
}

func (s *ObjectRuleV1Alpha_Mode) unmarshal(data []byte) error {
	str := strings.ToUpper(strings.TrimSpace(string(data)))
	if str == "" {
		*s = 0
		return nil
	}
	val, err := strconv.ParseInt(str, 10, 32)
	if err == nil {
		*s = ObjectRuleV1Alpha_Mode(val)
		return nil
	}
	val32, ok := ObjectRuleV1Alpha_Mode_value[str]
	if !ok {
		return errors.New("invalid state value")
	}
	*s = ObjectRuleV1Alpha_Mode(val32)
	return nil
}

func (s *ObjectRuleV1Alpha_Mode) UnmarshalJSON(data []byte) error {
	return s.unmarshal(data)
}

func (s *ObjectRuleV1Alpha_Mode) UnmarshalYAML(data []byte) error {
	return s.unmarshal(data)
}

func (s DeploymentTargetState) String() string {
	return DeploymentTargetState_name[int32(s)]
}

func (s SubjectRuleType) String() string {
	return SubjectRuleType_name[int32(s)]
}

func (s DeploymentState) String() string {
	return DeploymentState_name[int32(s)]
}

func (s ObjectRuleV1Alpha_Type) String() string {
	return ObjectRuleV1Alpha_Type_name[int32(s)]
}

func (s ObjectRuleV1Alpha_Mode) String() string {
	return ObjectRuleV1Alpha_Mode_name[int32(s)]
}
