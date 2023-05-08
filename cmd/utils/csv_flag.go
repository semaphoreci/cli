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

	if f == "" {
		return []string{}, nil
	}

	resultsWithWhiteSpace := strings.Split(f, ",")
	results := []string{}

	for _, r := range resultsWithWhiteSpace {
		results = append(results, strings.TrimSpace(r))
	}

	return results, err
}

func CSVArrayFlag(cmd *cobra.Command, flag string, trimSpace bool) (results [][]string, err error) {
	vals, err := cmd.Flags().GetStringArray(flag)
	if err != nil {
		return
	}
	for _, val := range vals {
		results = append(results, processCSVValue(val, trimSpace))
	}
	return
}

func processCSVValue(val string, trimSpace bool) (result []string) {
	parts := strings.Split(val, ",")
	for _, part := range parts {
		if trimSpace {
			part = strings.TrimSpace(part)
		}
		result = append(result, part)
	}

	return
}
