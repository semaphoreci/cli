package cmd

import (
	"net/http"
	"testing"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func Test__EditDashboard__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	dash := `{
		"metadata": {
			"name":"my-work",
			"id":"07e64c23-4bca-4f62-b76e-4a85aef24bbb",
			"title":"My Work",
			"create_time":"1536843668",
			"update_time":"1536843668"
		},
		"spec": {
			"widgets": [{
				"name": "Workflows",
				"type": "list",
				"filters":{
					"github_uid":"{{github_uid}}"
				}
			}]
		}
	}`

	received := false

	httpmock.RegisterResponder("GET", "https://org.semaphoretext.xyz/api/v1alpha/dashboards/my-work",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(200, dash), nil
		},
	)

	httpmock.RegisterResponder("PATCH", "https://org.semaphoretext.xyz/api/v1alpha/dashboards/07e64c23-4bca-4f62-b76e-4a85aef24bbb",
		func(req *http.Request) (*http.Response, error) {
			received = true

			return httpmock.NewStringResponse(200, dash), nil
		},
	)

	RootCmd.SetArgs([]string{"edit", "dashboards", "my-work"})
	RootCmd.Execute()

	if received == false {
		t.Error("Expected the API to receive GET and PATCH dashboard")
	}
}

func Test__EditSecret__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	dash := `{
		"metadata":{
			"name":"aaaaaaa",
			"id":"bb2ba294-d4b3-48bc-90a7-12dd56e9424b",
			"create_time":"1536673464",
			"update_time":"1536674946"
		},
		"data":{
			"env_vars":[{
				"name":"TEST",
				"value":"AAAA"
			}],
			"files":[{
				"path":"a.txt",
				"content":"W1twYWNrYWdlXV0KbmFtZSA9ICJzcGlyYWwtbWVtb3J5Igp2ZXJzaW9uID0gIjAuMS4wIgoK"
			}]
		}
	}`

	received := false

	httpmock.RegisterResponder("GET", "https://org.semaphoretext.xyz/api/v1beta/secrets/aaaaaaa",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(200, dash), nil
		},
	)

	httpmock.RegisterResponder("PATCH", "https://org.semaphoretext.xyz/api/v1beta/secrets/bb2ba294-d4b3-48bc-90a7-12dd56e9424b",
		func(req *http.Request) (*http.Response, error) {
			received = true

			return httpmock.NewStringResponse(200, dash), nil
		},
	)

	RootCmd.SetArgs([]string{"edit", "secrets", "aaaaaaa"})
	RootCmd.Execute()

	if received == false {
		t.Error("Expected the API to receive GET and PATCH secret")
	}
}
