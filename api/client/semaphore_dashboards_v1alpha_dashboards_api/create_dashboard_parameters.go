// Code generated by go-swagger; DO NOT EDIT.

package semaphore_dashboards_v1alpha_dashboards_api

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/renderedtext/sem/api/models"
)

// NewCreateDashboardParams creates a new CreateDashboardParams object
// with the default values initialized.
func NewCreateDashboardParams() *CreateDashboardParams {
	var ()
	return &CreateDashboardParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewCreateDashboardParamsWithTimeout creates a new CreateDashboardParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewCreateDashboardParamsWithTimeout(timeout time.Duration) *CreateDashboardParams {
	var ()
	return &CreateDashboardParams{

		timeout: timeout,
	}
}

// NewCreateDashboardParamsWithContext creates a new CreateDashboardParams object
// with the default values initialized, and the ability to set a context for a request
func NewCreateDashboardParamsWithContext(ctx context.Context) *CreateDashboardParams {
	var ()
	return &CreateDashboardParams{

		Context: ctx,
	}
}

// NewCreateDashboardParamsWithHTTPClient creates a new CreateDashboardParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewCreateDashboardParamsWithHTTPClient(client *http.Client) *CreateDashboardParams {
	var ()
	return &CreateDashboardParams{
		HTTPClient: client,
	}
}

/*CreateDashboardParams contains all the parameters to send to the API endpoint
for the create dashboard operation typically these are written to a http.Request
*/
type CreateDashboardParams struct {

	/*Body*/
	Body *models.SemaphoreDashboardsV1alphaDashboard

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the create dashboard params
func (o *CreateDashboardParams) WithTimeout(timeout time.Duration) *CreateDashboardParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the create dashboard params
func (o *CreateDashboardParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the create dashboard params
func (o *CreateDashboardParams) WithContext(ctx context.Context) *CreateDashboardParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the create dashboard params
func (o *CreateDashboardParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the create dashboard params
func (o *CreateDashboardParams) WithHTTPClient(client *http.Client) *CreateDashboardParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the create dashboard params
func (o *CreateDashboardParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the create dashboard params
func (o *CreateDashboardParams) WithBody(body *models.SemaphoreDashboardsV1alphaDashboard) *CreateDashboardParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the create dashboard params
func (o *CreateDashboardParams) SetBody(body *models.SemaphoreDashboardsV1alphaDashboard) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *CreateDashboardParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
