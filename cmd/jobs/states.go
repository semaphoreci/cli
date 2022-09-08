package jobs

import "strings"

func IsSelfHosted(machineType string) bool {
	return strings.HasPrefix(machineType, "s1-")
}
