package handler

import (
	"fmt"
)

type ProjectHandler struct {
}

func (h *ProjectHandler) Get(params GetParams) {
  fmt.Printf("Not implemented")
}

func (h *ProjectHandler) Describe(params DescribeParams) {
  fmt.Printf("Not implemented")
}

func (h *ProjectHandler) Create(params CreateParams) {
  fmt.Printf("Not implemented")
}

func (h *ProjectHandler) Delete(params DeleteParams) {
  fmt.Printf("Not implemented")
}
