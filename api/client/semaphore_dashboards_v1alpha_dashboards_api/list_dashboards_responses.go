// Code generated by go-swagger; DO NOT EDIT.

package semaphore_dashboards_v1alpha_dashboards_api

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/renderedtext/sem/api/models"
)

// ListDashboardsReader is a Reader for the ListDashboards structure.
type ListDashboardsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ListDashboardsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewListDashboardsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewListDashboardsOK creates a ListDashboardsOK with default headers values
func NewListDashboardsOK() *ListDashboardsOK {
	return &ListDashboardsOK{}
}

/*ListDashboardsOK handles this case with default header values.

(empty)
*/
type ListDashboardsOK struct {
	Payload *models.SemaphoreDashboardsV1alphaListDashboardsResponse
}

func (o *ListDashboardsOK) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha/dashboards][%d] listDashboardsOK  %+v", 200, o.Payload)
}

func (o *ListDashboardsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.SemaphoreDashboardsV1alphaListDashboardsResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
