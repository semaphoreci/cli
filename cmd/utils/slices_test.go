package utils

import (
	"testing"

	require "github.com/stretchr/testify/require"
)

func Test__Contains(t *testing.T) {
	require.False(t, Contains([]string{}, "a"))
	require.False(t, Contains([]string{"a", "b", "c"}, "d"))
	require.True(t, Contains([]string{"a", "b", "c"}, "a"))
}
