package client

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
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
	assert.True(t, received, "Expected the API to receive a GET request")
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
	assert.True(t, received, "Expected the API to receive a GET request")
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

	assert.True(t, received, "Expected the API to receive a GET request")
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

	assert.True(t, received, "Expected the API to receive a DELETE request")
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

	assert.True(t, received, "Expected the API to receive a POST request")
}

func Test__PostHeaders__InvalidURL(t *testing.T) {
	client := NewBaseClient("MYTOKEN", "org.semaphoretext.xyz", "testapi")

	response, statusCode, err := client.PostHeaders("^-^", []byte(""), map[string]string{})

	assert.Error(t, err)
	assert.Equal(t, 0, statusCode)
	assert.Equal(t, []byte(""), response)
}

func Test__Patch__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := false

	userAgent := "Test-User-Agent"
	UserAgent = userAgent
	client := NewBaseClient("MYTOKEN", "org.semaphoretext.xyz", "testapi")
	resourceBody := []byte("{}")

	httpmock.RegisterResponder("PATCH", "https://org.semaphoretext.xyz/api/testapi/kind/resource",
		func(req *http.Request) (*http.Response, error) {
			received = true
			assert.Equal(t, "Token MYTOKEN", req.Header.Get("Authorization"))
			assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
			assert.Equal(t, userAgent, req.Header.Get("User-Agent"))
			body, err := io.ReadAll(req.Body)
			assert.NoErrorf(t, err, "")
			assert.Equal(t, resourceBody, body)
			return httpmock.NewStringResponse(200, ""), nil
		},
	)

	response, statusCode, err := client.Patch("kind", "resource", resourceBody)

	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)
	assert.Equal(t, []byte(""), response)
	assert.True(t, received, "Expected the API to receive a PATCH request")
}

func Test__Patch__InvalidURL(t *testing.T) {
	client := NewBaseClient("MYTOKEN", "org.semaphoretext.xyz", "testapi")

	response, statusCode, err := client.Patch("^-^", "resource", []byte(""))

	assert.Error(t, err)
	assert.Equal(t, 0, statusCode)
	assert.Equal(t, []byte(""), response)
}

func Test__newfileUploadRequest__Success(t *testing.T) {
	// Lets create a temp file with the text "SOME CONTENT"
	f, err := os.CreateTemp("", "semtest")
	assert.NoError(t, err, "Could not create temporary file")
	defer os.Remove(f.Name())
	content := "SOME CONTENT"
	_, err = f.WriteString(content)
	assert.NoError(t, err, "Could not write to temporary file")

	// The parameters for newfileUploadRequest
	uri := "https://org.semaphoretext.xyz/api/upload"
	args := map[string]string{"key": "value"}
	fileArgName := "uploaded"
	path := f.Name()
	fileName := filepath.Base(path)

	// We explicitily define the type to make sure the return value matches
	var req *http.Request
	req, err = newfileUploadRequest(uri, args, fileArgName, path)

	// Lets check if request was correctly created
	assert.NoError(t, err)
	assert.Equal(t, "POST", req.Method)
	assert.Equal(t, "value", req.FormValue("key"))
	multiPartFile, multiPartHeader, err := req.FormFile("uploaded")
	assert.NoError(t, err, "Could not get the FormFile from the upload request: %s", err)
	assert.Equal(t, fileName, multiPartHeader.Filename)
	buffer := make([]byte, len(content)) // "SOME CONTENT"
	nbytes, err := multiPartFile.Read(buffer)
	assert.NoErrorf(t, err, "Could not read the multipart content: %s", err)
	assert.Equal(t, len(content), nbytes)
	assert.Equal(t, content, string(buffer))
}

func Test__newfileUploadRequest__FailToOpenFile(t *testing.T) {
	// Lets create a temp file, keep the path then delete it
	// to make sure we have and invalid path
	f, err := os.CreateTemp("", "semtest")
	assert.NoError(t, err, "Could not create temporary file")
	// Just in case test crashes prematurely
	defer os.Remove(f.Name())

	// The parameters for newfileUploadRequest
	uri := "https://org.semaphoretext.xyz/api/upload"
	args := map[string]string{"key": "value"}
	fileArgName := "uploaded"
	path := f.Name() // lets keep the path to the temp file

	// Delete the file
	assert.NoError(t, os.Remove(f.Name()))

	req, err := newfileUploadRequest(uri, args, fileArgName, path)

	fmt.Printf("Error: %s", err)
	assert.Error(t, err)
	var result *http.Request = nil
	assert.Equal(t, result, req)
}
