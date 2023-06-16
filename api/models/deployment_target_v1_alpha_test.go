package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeploymentTargetsToJson(t *testing.T) {
	dt := DeploymentTargetV1Alpha{
		ApiVersion: "v1alpha",
		Kind:       DeploymentTargetKindV1Alpha,
		DeploymentTargetMetadataV1Alpha: DeploymentTargetMetadataV1Alpha{
			Id:             "1234-5678-id",
			Name:           "dt-name",
			OrganizationId: "org-id",
			ProjectId:      "prj-id",
			Description:    "dt-description",
			Url:            "www.semaphore.xyz",
		},
		DeploymentTargetSpecV1Alpha: DeploymentTargetSpecV1Alpha{
			Active:             false,
			BookmarkParameter1: "book1",
			SubjectRules: []*SubjectRuleV1Alpha{{
				Type:      "USER",
				SubjectId: "00000000-0000-0000-0000-000000000000",
			}},
			ObjectRules: []*ObjectRuleV1Alpha{
				{Type: "BRANCH", MatchMode: "EXACT", Pattern: ".*main.*"},
			},
			State: "CORDONED",
		},
	}

	received_json, err := dt.ToJson()
	assert.Nil(t, err)

	expected_json := `{"id":"1234-5678-id","name":"dt-name","project_id":"prj-id","organization_id":"org-id","description":"dt-description","url":"www.semaphore.xyz","state":"CORDONED","state_message":"","subject_rules":[{"type":"USER","subject_id":"00000000-0000-0000-0000-000000000000"}],"object_rules":[{"type":"BRANCH","match_mode":"EXACT","pattern":".*main.*"}],"active":false,"bookmark_parameter1":"book1","bookmark_parameter2":"","bookmark_parameter3":""}`
	assert.Equal(t, expected_json, string(received_json))
}

func TestDeploymentTargetsToYaml(t *testing.T) {
	dt := DeploymentTargetV1Alpha{
		ApiVersion: "v1alpha",
		Kind:       DeploymentTargetKindV1Alpha,
		DeploymentTargetMetadataV1Alpha: DeploymentTargetMetadataV1Alpha{
			Id:             "1234-5678-id",
			Name:           "dt-name",
			OrganizationId: "org-id",
			ProjectId:      "prj-id",
			Description:    "dt-description",
			Url:            "www.semaphore.xyz",
		},
		DeploymentTargetSpecV1Alpha: DeploymentTargetSpecV1Alpha{
			Active:             true,
			BookmarkParameter1: "book1",
			SubjectRules: []*SubjectRuleV1Alpha{{
				Type:      "USER",
				SubjectId: "subjId1",
			}},
			ObjectRules: []*ObjectRuleV1Alpha{
				{Type: "BRANCH", MatchMode: "EXACT", Pattern: ".*main.*"},
			},
			State: "USABLE",
		},
	}
	received_yaml, err := dt.ToYaml()
	assert.Nil(t, err)

	expected_yaml := `apiVersion: v1alpha
kind: DeploymentTarget
metadata:
  id: 1234-5678-id
  name: dt-name
  project_id: prj-id
  organization_id: org-id
  description: dt-description
  url: www.semaphore.xyz
spec:
  state: USABLE
  state_message: ""
  subject_rules:
  - type: USER
    subject_id: subjId1
  object_rules:
  - type: BRANCH
    match_mode: EXACT
    pattern: .*main.*
  active: true
  bookmark_parameter1: book1
  bookmark_parameter2: ""
  bookmark_parameter3: ""
`
	assert.Equal(t, expected_yaml, string(received_yaml))
}

func TestDeploymentTargetsFromYaml(t *testing.T) {
	content := `apiVersion: v1alpha
kind: DeploymentTarget
metadata:
  id: 1234-5678-id
  name: dt-name
  organization_id: org-id
  project_id: prj-id
  url: www.semaphore.xyz
  description: dt-description
spec:
  active: true
  bookmark_parameter1: book1
`
	dt, err := NewDeploymentTargetV1AlphaFromYaml([]byte(content))
	assert.Nil(t, err)

	assert.Equal(t, dt.Id, "1234-5678-id")
	assert.Equal(t, dt.Name, "dt-name")
	assert.Equal(t, dt.ProjectId, "prj-id")
	assert.Equal(t, dt.OrganizationId, "org-id")
	assert.Equal(t, dt.Url, "www.semaphore.xyz")
	assert.Equal(t, dt.Description, "dt-description")
	assert.True(t, dt.Active)
	assert.Equal(t, dt.BookmarkParameter1, "book1")
}

func TestDeploymentTargetsFromJSON(t *testing.T) {
	content := `{"id":"1234-5678-id","name":"dt-name","project_id":"prj-id","organization_id":"org-id","description":"dt-description","url":"www.semaphore.xyz","state":"USABLE","subject_rules":[{"type":"USER","subject_id":"00000000-0000-0000-0000-000000000000"}],"object_rules":[{"type":"BRANCH","match_mode":"REGEX","pattern":".*main.*"}],"active":true,"bookmark_parameter1":"book1","env_vars":[{"name":"Var1","value":"Val1"}],"files":[{"path":"/etc/config.yml","content":"abcdefgh"}]}`
	dt, err := NewDeploymentTargetV1AlphaFromJson([]byte(content))
	assert.Nil(t, err)

	assert.Equal(t, dt.Id, "1234-5678-id")
	assert.Equal(t, dt.Name, "dt-name")
	assert.Equal(t, dt.ProjectId, "prj-id")
	assert.Equal(t, dt.OrganizationId, "org-id")
	assert.Equal(t, dt.Url, "www.semaphore.xyz")
	assert.Equal(t, dt.Description, "dt-description")
	assert.True(t, dt.Active)
	assert.Equal(t, dt.BookmarkParameter1, "book1")
}
