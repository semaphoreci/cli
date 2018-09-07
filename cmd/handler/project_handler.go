package handler

import (
	"encoding/json"
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/renderedtext/sem/client"
	"github.com/renderedtext/sem/cmd/utils"
)

type ProjectHandler struct {
}

func (h *ProjectHandler) Get(params GetParams) {
	c := client.FromConfig()

	body, status, err := c.List("projects")

	utils.Check(err)

	if status != 200 {
		utils.Fail(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	var secrets []map[string]interface{}

	json.Unmarshal([]byte(body), &secrets)

	fmt.Println("NAME")

	for _, secret := range secrets {
		fmt.Println(secret["metadata"].(map[string]interface{})["name"])
	}
}

func (h *ProjectHandler) Describe(params DescribeParams) {
	c := client.FromConfig()

	body, status, err := c.Get("projects", params.Name)

	utils.Check(err)

	if status != 200 {
		utils.Fail(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	j, _ := yaml.JSONToYAML(body)

	fmt.Println(string(j))
}

func (h *ProjectHandler) Create(params CreateParams) {
	c := client.FromConfig()
	c.SetApiVersion(params.ApiVersion)

	body, status, err := c.Post("projects", params.Resource)

	utils.Check(err)

	if status != 200 {
		utils.Fail(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	j, _ := yaml.JSONToYAML(body)

	fmt.Println(string(j))
}

func (h *ProjectHandler) Edit(params EditParams) {
	fmt.Printf("Unsupported Action")
}

func (h *ProjectHandler) Delete(params DeleteParams) {
	c := client.FromConfig()

	body, status, err := c.Delete("projects", params.Name)

	utils.Check(err)

	if status != 200 {
		utils.Fail(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	fmt.Printf("project \"%s\" deleted\n", params.Name)
}
