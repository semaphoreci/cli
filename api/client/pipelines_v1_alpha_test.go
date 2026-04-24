package client

import (
	"net/http"
	"testing"
)

func TestHasNextPage(t *testing.T) {
	tests := []struct {
		name     string
		headers  http.Header
		expected bool
	}{
		{
			name:     "nil headers",
			headers:  nil,
			expected: false,
		},
		{
			name:     "no Link header",
			headers:  http.Header{},
			expected: false,
		},
		{
			name: "Link header with rel=next",
			headers: http.Header{
				"Link": []string{
					`<http://example.com/api/v1alpha/pipelines?page=2>; rel="next", <http://example.com/api/v1alpha/pipelines?page=1>; rel="first"`,
				},
			},
			expected: true,
		},
		{
			name: "Link header without rel=next",
			headers: http.Header{
				"Link": []string{
					`<http://example.com/api/v1alpha/pipelines?page=1>; rel="first", <http://example.com/api/v1alpha/pipelines?page=5>; rel="last"`,
				},
			},
			expected: false,
		},
		{
			name: "Link header with only rel=last",
			headers: http.Header{
				"Link": []string{
					`<http://example.com/api/v1alpha/pipelines?page=5>; rel="last"`,
				},
			},
			expected: false,
		},
		{
			name: "multiple Link header values",
			headers: http.Header{
				"Link": []string{
					`<http://example.com/api/v1alpha/pipelines?page=1>; rel="first"`,
					`<http://example.com/api/v1alpha/pipelines?page=3>; rel="next"`,
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasNextPage(tt.headers)
			if result != tt.expected {
				t.Errorf("hasNextPage() = %v, want %v", result, tt.expected)
			}
		})
	}
}
