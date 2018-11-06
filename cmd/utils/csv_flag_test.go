package utils

import (
	"testing"

	"github.com/spf13/cobra"
	assert "github.com/stretchr/testify/assert"
)

func Test__CSVFlag__SplitsString(t *testing.T) {
	cmd := &cobra.Command{}

	cmd.Flags().String("test", "", "")
	cmd.SetArgs([]string{"--test", "a,b,c,d"})
	cmd.Execute()

	values, err := CSVFlag(cmd, "test")

	assert.Nil(t, err)
	assert.Equal(t, values, []string{"a", "b", "c", "d"})
}

func Test__CSVFlag__RemovesWhitespace(t *testing.T) {
	cmd := &cobra.Command{}

	cmd.Flags().String("test", "", "")
	cmd.SetArgs([]string{"--test", "a, b, c, d"})
	cmd.Execute()

	values, err := CSVFlag(cmd, "test")

	assert.Nil(t, err)
	assert.Equal(t, values, []string{"a", "b", "c", "d"})
}

func Test__CSVFlag__EmptyStringHasNoElements(t *testing.T) {
	cmd := &cobra.Command{}

	cmd.Flags().String("test", "", "")
	cmd.Execute()

	values, err := CSVFlag(cmd, "test")

	assert.Nil(t, err)
	assert.Equal(t, values, []string{})
}
