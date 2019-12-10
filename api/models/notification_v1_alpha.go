package models

import (
	"encoding/json"
	"fmt"

	yaml "gopkg.in/yaml.v2"
)

type NotificationV1AlphaMetadata struct {
	Name       string      `json:"name,omitempty" yaml:"name,omitempty"`
	Id         string      `json:"id,omitempty" yaml:"id,omitempty"`
	CreateTime json.Number `json:"create_time,omitempty,string" yaml:"create_time,omitempty"`
	UpdateTime json.Number `json:"update_time,omitempty,string" yaml:"update_time,omitempty"`
}

type NotificationV1AlphaSpecRuleFilter struct {
	Projects  []string `json:"projects,omitempty" yaml:"projects,omitempty"`
	Branches  []string `json:"branches,omitempty" yaml:"branches,omitempty"`
	Pipelines []string `json:"pipelines,omitempty" yaml:"pipelines,omitempty"`
	Blocks    []string `json:"blocks,omitempty" yaml:"blocks,omitempty"`
	States    []string `json:"states,omitempty" yaml:"states,omitempty"`
	Results   []string `json:"results,omitempty" yaml:"results,omitempty"`
}

type NotificationV1AlphaSpecRuleNotify struct {
	Slack struct {
		Endpoint string   `json:"endpoint,omitempty" yaml:"endpoint,omitempty"`
		Channels []string `json:"channels,omitempty" yaml:"channels,omitempty"`
		Message  string   `json:"message,omitempty" yaml:"message,omitempty"`
	} `json:"slack,omitempty" yaml:"slack,omitempty"`

	Email struct {
		Subject string   `json:"subject,omitempty" yaml:"subject,omitempty"`
		CC      []string `json:"cc,omitempty" yaml:"cc,omitempty"`
		BCC     []string `json:"bcc,omitempty" yaml:"bcc,omitempty"`
		Content string   `json:"content,omitempty" yaml:"content,omitempty"`
	} `json:"email,omitempty" yaml:"email,omitempty"`

	Webhook struct {
		Endpoint  string   `json:"endpoint,omitempty" yaml:"endpoint,omitempty"`
		Timeout   int32    `json:"timeout,omitempty" yaml:"timeout,omitempty"`
		Action    string    `json:"action,omitempty" yaml:"action,omitempty"`
		Retries   int32    `json:"retries,omitempty" yaml:"retries,omitempty"`
	} `json:"webhook,omitempty" yaml:"webhook,omitempty"`
}

type NotificationV1AlphaSpecRule struct {
	Name   string                            `json:"name,omitempty" yaml:"name,omitempty"`
	Filter NotificationV1AlphaSpecRuleFilter `json:"filter,omitempty" yaml:"filter,omitempty"`
	Notify NotificationV1AlphaSpecRuleNotify `json:"notify,omitempty" yaml:"notify,omitempty"`
}

type NotificationV1AlphaSpec struct {
	Rules []NotificationV1AlphaSpecRule `json:"rules,omitempty" yaml:"rules,omitempty"`
}

type NotificationV1AlphaStatusFailure struct {
	Time    json.Number `json:"time,omitempty,string" yaml:"time,omitempty"`
	Message string      `json:"message,omitempty,string" yaml:"message,omitempty"`
}

type NotificationV1AlphaStatus struct {
	Failures []NotificationV1AlphaStatusFailure `json:"failures,omitempty" yaml:"failures,omitempty"`
}

type NotificationV1Alpha struct {
	ApiVersion string `json:"apiVersion,omitempty" yaml:"apiVersion"`
	Kind       string `json:"kind,omitempty" yaml:"kind"`

	Metadata NotificationV1AlphaMetadata `json:"metadata" yaml:"metadata"`
	Spec     NotificationV1AlphaSpec     `json:"spec" yaml:"spec"`
	Status   NotificationV1AlphaStatus   `json:"status" yaml:"status"`
}

func NewNotificationV1Alpha(name string) *NotificationV1Alpha {
	n := NotificationV1Alpha{}

	n.setApiVersionAndKind()
	n.Metadata.Name = name

	return &n
}

func NewNotificationV1AlphaFromJson(data []byte) (*NotificationV1Alpha, error) {
	n := NotificationV1Alpha{}

	err := json.Unmarshal(data, &n)

	if err != nil {
		return nil, err
	}

	n.setApiVersionAndKind()

	return &n, nil
}

func NewNotificationV1AlphaFromYaml(data []byte) (*NotificationV1Alpha, error) {
	n := NotificationV1Alpha{}

	err := yaml.UnmarshalStrict(data, &n)

	if err != nil {
		return nil, err
	}

	n.setApiVersionAndKind()

	return &n, nil
}

func (n *NotificationV1Alpha) setApiVersionAndKind() {
	n.ApiVersion = "v1alpha"
	n.Kind = "Notification"
}

func (n *NotificationV1Alpha) ObjectName() string {
	return fmt.Sprintf("Notifications/%s", n.Metadata.Name)
}

func (n *NotificationV1Alpha) ToJson() ([]byte, error) {
	return json.Marshal(n)
}

func (n *NotificationV1Alpha) ToYaml() ([]byte, error) {
	return yaml.Marshal(n)
}
