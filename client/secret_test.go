package client

import "testing"

func TestInitSecretFromYaml__ValidYaml(t *testing.T) {
	yaml := `
apiVersion: "v1beta"
kind: Secret
metadata:
  name: test
data:
  env_vars:
    - name: A
      value: B

  files:
    - path: "a.txt"
      content: "31312312dfadsfa3412323"
`
	secret, err := InitSecretFromYaml([]byte(yaml))

	if err != nil {
		t.Errorf("Unexpected error while unmarshaling secret '%s'", err)
	}

	if secret.Metadata.Name != "test" {
		t.Errorf("Name is incorrect, got: %s, want: %s.", secret.Metadata.Name, "test")
	}

	if secret.Data.EnvVars[0].Name != "A" {
		t.Errorf("Invalid env var, got: '%s', want: %s", secret.Data.EnvVars[0].Name, "A")
	}

	if secret.Data.EnvVars[0].Value != "B" {
		t.Errorf("Invalid env var, got: '%s', want: %s", secret.Data.EnvVars[0].Name, "B")
	}

	if secret.Data.Files[0].Path != "a.txt" {
		t.Errorf("Invalid file path, got: '%s', want: %s", secret.Data.Files[0].Path, "a.txt")
	}

	if secret.Data.Files[0].Content != "31312312dfadsfa3412323" {
		t.Errorf("Invalid file content, got: '%s', want: %s", secret.Data.Files[0].Content, "31312312dfadsfa3412323")
	}
}

func TestInitSecretFromYaml__InvalidYaml(t *testing.T) {
	yaml := `
data:
  env_vars:
    - name: A
      value: B
  tests:
    - path: "a.txt"
      content: "31312312dfadsfa3412323"
`
	_, err := InitSecretFromYaml([]byte(yaml))

	if err == nil {
		t.Errorf("Expected error but got nil")
	}
}
