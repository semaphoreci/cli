// Code generated by go-swagger; DO NOT EDIT.

package semaphore_secrets_v1beta_secrets_api

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// New creates a new semaphore secrets v1beta secrets api API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) *Client {
	return &Client{transport: transport, formats: formats}
}

/*
Client for semaphore secrets v1beta secrets api API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

/*
CreateSecret create secret API
*/
func (a *Client) CreateSecret(params *CreateSecretParams) (*CreateSecretOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewCreateSecretParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "CreateSecret",
		Method:             "POST",
		PathPattern:        "/api/v1beta/secrets",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &CreateSecretReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*CreateSecretOK), nil

}

/*
DeleteSecret delete secret API
*/
func (a *Client) DeleteSecret(params *DeleteSecretParams) (*DeleteSecretOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewDeleteSecretParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "DeleteSecret",
		Method:             "DELETE",
		PathPattern:        "/api/v1beta/secrets/{secret_id_or_name}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &DeleteSecretReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*DeleteSecretOK), nil

}

/*
GetSecret get secret API
*/
func (a *Client) GetSecret(params *GetSecretParams) (*GetSecretOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetSecretParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "GetSecret",
		Method:             "GET",
		PathPattern:        "/api/v1beta/secrets/{secret_id_or_name}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &GetSecretReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*GetSecretOK), nil

}

/*
ListSecrets list secrets API
*/
func (a *Client) ListSecrets(params *ListSecretsParams) (*ListSecretsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListSecretsParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "ListSecrets",
		Method:             "GET",
		PathPattern:        "/api/v1beta/secrets",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &ListSecretsReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*ListSecretsOK), nil

}

/*
UpdateSecret update secret API
*/
func (a *Client) UpdateSecret(params *UpdateSecretParams) (*UpdateSecretOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewUpdateSecretParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "UpdateSecret",
		Method:             "PATCH",
		PathPattern:        "/api/v1beta/secrets/{secret_id_or_name}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &UpdateSecretReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*UpdateSecretOK), nil

}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
