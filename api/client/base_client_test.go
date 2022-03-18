package client

import (
	"io"
	"net/http"
	"net/url"
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
		t.Error("Expected the API to receive a GET request")
	}
}

func Test__Get__InvalidURL(t *testing.T) {
	client := NewBaseClient("MYTOKEN", "org.semaphoretext.xyz", "testapi")

	response, statusCode, err := client.Get("kind", "^-^")

	assert.Error(t, err)
	assert.Equal(t, 0, statusCode)
	assert.Equal(t, []byte(""), response)
}

func Test__List__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := false

	userAgent := "Test-User-Agent"
	UserAgent = userAgent

	httpmock.RegisterResponder("GET", "https://org.semaphoretext.xyz/api/testapi/kind",
		func(req *http.Request) (*http.Response, error) {
			received = true
			assert.Equal(t, "Token MYTOKEN", req.Header.Get("Authorization"))
			assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
			assert.Equal(t, userAgent, req.Header.Get("User-Agent"))

			return httpmock.NewStringResponse(200, ""), nil
		},
	)

	client := NewBaseClient("MYTOKEN", "org.semaphoretext.xyz", "testapi")

	response, statusCode, err := client.List("kind")

	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)
	assert.Equal(t, []byte(""), response)

	if received == false {
		t.Error("Expected the API to receive a GET request")
	}
}

func Test__List__InvalidURL(t *testing.T) {
	client := NewBaseClient("MYTOKEN", "org.semaphoretext.xyz", "testapi")

	response, statusCode, err := client.List("^-^")

	assert.Error(t, err)
	assert.Equal(t, 0, statusCode)
	assert.Equal(t, []byte(""), response)
}

func Test__ListWithParams__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := false

	userAgent := "Test-User-Agent"
	UserAgent = userAgent

	values := url.Values{}
	values.Add("some", "value")

	httpmock.RegisterResponder("GET", "https://org.semaphoretext.xyz/api/testapi/kind",
		func(req *http.Request) (*http.Response, error) {
			received = true
			assert.Equal(t, "Token MYTOKEN", req.Header.Get("Authorization"))
			assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
			assert.Equal(t, userAgent, req.Header.Get("User-Agent"))
			assert.Equal(t, "value", req.URL.Query().Get("some"))
			return httpmock.NewStringResponse(200, ""), nil
		},
	)

	client := NewBaseClient("MYTOKEN", "org.semaphoretext.xyz", "testapi")

	response, statusCode, err := client.ListWithParams("kind", values)

	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)
	assert.Equal(t, []byte(""), response)

	if received == false {
		t.Error("Expected the API to receive a GET request")
	}
}

func Test__ListWithParams__InvalidURL(t *testing.T) {
	client := NewBaseClient("MYTOKEN", "org.semaphoretext.xyz", "testapi")

	response, statusCode, err := client.ListWithParams("^-^", url.Values{})

	assert.Error(t, err)
	assert.Equal(t, 0, statusCode)
	assert.Equal(t, []byte(""), response)
}

func Test__Delete__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := false

	userAgent := "Test-User-Agent"
	UserAgent = userAgent

	values := url.Values{}
	values.Add("some", "value")

	httpmock.RegisterResponder("DELETE", "https://org.semaphoretext.xyz/api/testapi/kind/resource",
		func(req *http.Request) (*http.Response, error) {
			received = true
			assert.Equal(t, "Token MYTOKEN", req.Header.Get("Authorization"))
			assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
			assert.Equal(t, userAgent, req.Header.Get("User-Agent"))
			return httpmock.NewStringResponse(200, ""), nil
		},
	)

	client := NewBaseClient("MYTOKEN", "org.semaphoretext.xyz", "testapi")

	response, statusCode, err := client.Delete("kind", "resource")

	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)
	assert.Equal(t, []byte(""), response)

	if received == false {
		t.Error("Expected the API to receive a DELETE request")
	}
}

func Test__Delete__InvalidURL(t *testing.T) {
	client := NewBaseClient("MYTOKEN", "org.semaphoretext.xyz", "testapi")

	response, statusCode, err := client.Delete("kind", "^-^")

	assert.Error(t, err)
	assert.Equal(t, 0, statusCode)
	assert.Equal(t, []byte(""), response)
}

func Test__PostHeaders__Response201(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := false

	userAgent := "Test-User-Agent"
	UserAgent = userAgent
	headers := map[string]string{"key": "value"}
	client := NewBaseClient("MYTOKEN", "org.semaphoretext.xyz", "testapi")
	resource := []byte("resource")

	httpmock.RegisterResponder("POST", "https://org.semaphoretext.xyz/api/testapi/kind",
		func(req *http.Request) (*http.Response, error) {
			received = true
			assert.Equal(t, "Token MYTOKEN", req.Header.Get("Authorization"))
			assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
			assert.Equal(t, userAgent, req.Header.Get("User-Agent"))
			assert.Equal(t, "value", req.Header.Get("key"))

			body, err := io.ReadAll(req.Body)
			assert.NoError(t, err)
			assert.Equal(t, "resource", string(body))
			return httpmock.NewStringResponse(201, ""), nil
		},
	)

	response, statusCode, err := client.PostHeaders("kind", resource, headers)

	assert.NoError(t, err)
	assert.Equal(t, 201, statusCode)
	assert.Equal(t, []byte(""), response)

	if received == false {
		t.Error("Expected the API to receive a POST request")
	}
}

func Test__PostHeaders__InvalidURL(t *testing.T) {
	client := NewBaseClient("MYTOKEN", "org.semaphoretext.xyz", "testapi")

	response, statusCode, err := client.PostHeaders("^-^", []byte(""), map[string]string{})

	assert.Error(t, err)
	assert.Equal(t, 0, statusCode)
	assert.Equal(t, []byte(""), response)
}
