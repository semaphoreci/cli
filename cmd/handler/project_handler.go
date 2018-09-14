package handler

import (
	"encoding/json"
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/semaphoreci/cli/client"
	"github.com/semaphoreci/cli/cmd/utils"
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

	fmt.Print(string(j))
}

func (h *ProjectHandler) Create(params CreateParams) {
	p, err := client.InitProjectFromYaml(params.Resource)

	utils.Check(err)

	err = p.Create()

	utils.Check(err)

	fmt.Printf("project \"%s\" created\n", p.Metadata.Name)
}

func (h *ProjectHandler) Apply(params ApplyParams) {
	fmt.Printf("error: Unsupported Action\n")
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
