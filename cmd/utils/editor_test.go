package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Test__newEditorCommand__PathWithSpacesIsSingleArgument(t *testing.T) {
	tmpDir := t.TempDir()

	logPath := filepath.Join(tmpDir, "editor-args.log")
	scriptPath := filepath.Join(tmpDir, "capture-args.sh")
	filePath := filepath.Join(tmpDir, "Notifications-my notification.yml")

	script := `#!/bin/sh
log_path="$1"
shift
printf '%s\n' "$#" > "$log_path"
for arg in "$@"; do
  printf '%s\n' "$arg" >> "$log_path"
done
`

	if err := os.WriteFile(scriptPath, []byte(script), 0o755); err != nil {
		t.Fatalf("failed to create capture script: %v", err)
	}

	editor := scriptPath + " " + logPath
	cmd := newEditorCommand(editor, filePath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("editor command failed: %v (output: %s)", err, string(output))
	}

	logData, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("failed to read captured args: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(logData)), "\n")
	if len(lines) < 2 {
		t.Fatalf("expected at least 2 lines in capture output, got %d", len(lines))
	}

	if lines[0] != "1" {
		t.Fatalf("expected editor script to receive exactly one file path argument, got %q", lines[0])
	}

	if lines[1] != filePath {
		t.Fatalf("expected file path %q, got %q", filePath, lines[1])
	}
}
