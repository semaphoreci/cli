package handler

import (
	"fmt"
  "errors"
)

type GetParams struct {
}

type DescribeParams struct {
  Name string
}

type DeleteParams struct {
  Name string
}

type CreateParams struct {
  ApiVersion string
  Resource []byte
}

type Handler interface {
  Get(GetParams)
  Describe(DescribeParams)
  Create(CreateParams)
  Delete(DeleteParams)
}

func FindHandler(resource_kind string) (Handler, error) {
  switch resource_kind {
    case "secret", "secrets":
      return new(SecretHandler), nil
    case "project", "projects":
      return new(ProjectHandler), nil
    default:
      return nil, errors.New(fmt.Sprintf("Unknown resource kind %s.", resource_kind))
  }
}
