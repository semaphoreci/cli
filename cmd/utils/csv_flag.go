package utils

import (
	"strings"

	"github.com/spf13/cobra"
)

func CSVFlag(cmd *cobra.Command, flag string) ([]string, error) {
	f, err := cmd.Flags().GetString(flag)

	if err != nil {
		return []string{}, err
	}

	resultsWithWhiteSpace := strings.Split(f, ",")
	results := []string{}

	for _, r := range resultsWithWhiteSpace {
		results = append(results, strings.TrimSpace(r))
	}

	return results, err
}
