package handler

import (
	"errors"
	"fmt"
	"strings"
)

type GetParams struct {
}

type DescribeParams struct {
	Name string
}

type EditParams struct {
	Name string
}

type DeleteParams struct {
	Name string
}

type CreateParams struct {
	ApiVersion string
	Resource   []byte
}

type Handler interface {
	Get(GetParams)
	Describe(DescribeParams)
	Create(CreateParams)
	Edit(EditParams)
	Delete(DeleteParams)
}

func FindHandler(resource_kind string) (Handler, error) {
	switch strings.ToLower(resource_kind) {
	case "secret", "secrets":
		return new(SecretHandler), nil
	case "project", "projects", "prj":
		return new(ProjectHandler), nil
	default:
		return nil, errors.New(fmt.Sprintf("Unknown resource kind %s.", resource_kind))
	}
}
