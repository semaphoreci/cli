package client

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func Test__Get__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := false

	userAgent := "Test-User-Agent"
	UserAgent = userAgent

	httpmock.RegisterResponder("GET", "https://org.semaphoretext.xyz/api/testapi/kind/resource",
		func(req *http.Request) (*http.Response, error) {
			received = true
			assert.Equal(t, "Token MYTOKEN", req.Header.Get("Authorization"))
			assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
			assert.Equal(t, userAgent, req.Header.Get("User-Agent"))

			return httpmock.NewStringResponse(200, ""), nil
		},
	)

	client := NewBaseClient("MYTOKEN", "org.semaphoretext.xyz", "testapi")

	response, statusCode, err := client.Get("kind", "resource")

	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)
	assert.Equal(t, []byte(""), response)

	if received == false {
		t.Error("Expected the API to receive GET request")
	}
}

func Test__Get__InvalidURL(t *testing.T) {
	client := NewBaseClient("MYTOKEN", "org.semaphoretext.xyz", "testapi")

	response, statusCode, err := client.Get("kind", "^-^")

	assert.Error(t, err)
	assert.Equal(t, 0, statusCode)
	assert.Equal(t, []byte(""), response)
}
