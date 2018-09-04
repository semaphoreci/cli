package client

import "testing"

func TestInitProjectFromYaml(t *testing.T) {
	yaml := `apiVersion: v1alpha
kind: Project
metadata:
  name: test
spec:
  repository:
    url: "git@github.com:/renderedtext/sem.git"
`
	project, _ := InitProjectFromYaml([]byte(yaml))

	if project.Metadata.Name != "test" {
		t.Errorf("Name is incorrect, got: %s, want: %s.", project.Metadata.Name, "test")
	}

	if project.Spec.Repository.Url != "git@github.com:/renderedtext/sem.git" {
		t.Errorf("Repo Url is incorrect, got: %s, want: %s.", project.Spec.Repository.Url, "git@github.com:/renderedtext/sem.git")
	}
}

func TestToYaml(t *testing.T) {
	project := InitProject("test", "github.com:/renderedtext/sem.git")

	json_body, _ := project.ToJson()

	expected_json_body := `{"metadata":{"name":"test"},"spec":{"repository":{"url":"github.com:/renderedtext/sem.git"}}}`

	if string(json_body) != expected_json_body {
		t.Errorf("JSON body is incorrect, got: %s, want: %s.", json_body, expected_json_body)
	}
}
