package generators

import "testing"

func TestConstructProjectName(t *testing.T) {
	name, err := ConstructProjectName("git@github.com:/renderedtext/test.git")

	if name != "test" {
		t.Errorf("Name is incorrect, got: %s, want: %s.", name, "test")
	}

	name, err = ConstructProjectName("git@github.com:/renderedtext/test")

	if name != "test" {
		t.Errorf("Name is incorrect, got: %s, want: %s.", name, "test")
	}

	name, err = ConstructProjectName("github.com/renderedtext/test")

	if err == nil {
		t.Errorf("Expected error for unrecognized format.")
	}
}
