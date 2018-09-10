package handler

import (
	"fmt"
	"time"
)

func RelativeAgeForHumans(timestamp int64) string {
	seconds := currentTimestamp() - timestamp

	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	}

	minutes := seconds / 60

	if minutes < 60 {
		return fmt.Sprintf("%dm", minutes)
	}

	hours := minutes / 60

	if hours < 24 {
		return fmt.Sprintf("%dh", hours)
	}

	days := hours / 24

	return fmt.Sprintf("%dd", days)
}

func currentTimestamp() int64 {
	return time.Now().UnixNano() / 1e9
}
