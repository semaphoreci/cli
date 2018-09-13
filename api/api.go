package api

import (
	"fmt"
	"net/http"

	"github.com/go-openapi/strfmt"

	httptransport "github.com/go-openapi/runtime/client"
	apiclient "github.com/renderedtext/sem/api/client"
	"github.com/renderedtext/sem/config"
)

func DefaultClient() *apiclient.Semaphore {
	host := config.GetHost()
	transport := httptransport.New(host, "", []string{"https"})
	transport.Transport = newRoundTripper()

	return apiclient.New(transport, strfmt.Default)
}

type roundtripper struct {
	underlyingTransport http.RoundTripper
}

func newRoundTripper() *roundtripper {
	return &roundtripper{underlyingTransport: http.DefaultTransport}
}

func (rt *roundtripper) RoundTrip(r *http.Request) (*http.Response, error) {
	// fmt.Println("Sending request")

	r.Header.Add("Authorization", fmt.Sprintf("Token %s", config.GetAuth()))
	res, err := rt.underlyingTransport.RoundTrip(r)

	// fmt.Println("Received response")

	return res, err
}
