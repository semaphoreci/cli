package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/semaphoreci/cli/config"
)

const editedContentTemplate = `# Editing %s.
# When you close the editor, the content will be updated on %s.

%s
`

func EditYamlInEditor(objectName string, content string) (string, error) {
	content_with_comment := fmt.Sprintf(
		editedContentTemplate,
		objectName,
		config.GetHost(),
		content)

	dir, err := ioutil.TempDir("", "sem-cli-session")

	if err != nil {
		return "", fmt.Errorf("Failed to open local temp file for editing '%s'", err)
	}

	defer os.RemoveAll(dir) // clean up

	// remove '/' from filename
	filename := strings.Replace(fmt.Sprintf("%s.yml", objectName), "/", "-", -1)
	// #nosec
	tmpfile, err := os.Create(filepath.Join(dir, filename))

	if err != nil {
		return "", fmt.Errorf("Failed to open local temp file for editing '%s'", err)
	}

	if _, err := tmpfile.Write([]byte(content_with_comment)); err != nil {
		return "", fmt.Errorf("Failed to open local temp file for editing '%s'", err)
	}

	if err := tmpfile.Close(); err != nil {
		return "", fmt.Errorf("Failed to open local temp file for editing '%s'", err)
	}

	editor := config.GetEditor()

	cmd := exec.Command("sh", "-c", fmt.Sprintf("%s %s", editor, tmpfile.Name()))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Start()

	if err != nil {
		return "", fmt.Errorf("Failed to start editor '%s'", err)
	}

	err = cmd.Wait()

	if err != nil {
		return "", fmt.Errorf("Editor closed with with '%s'", err)
	}

	editedContent, err := ioutil.ReadFile(tmpfile.Name())

	if err != nil {
		return "", fmt.Errorf("Failed to read the content of the edited object '%s'", err)
	}

	return string(editedContent), nil
}
