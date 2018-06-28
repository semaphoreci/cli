package handler

import (
	"encoding/json"
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/renderedtext/sem/client"
)

type ProjectHandler struct {
}

func (h *ProjectHandler) Get(params GetParams) {
	c := client.FromConfig()

	body, _, _ := c.List("projects")

	var secrets []map[string]interface{}

	json.Unmarshal([]byte(body), &secrets)

	fmt.Println("NAME")

	for _, secret := range secrets {
		fmt.Println(secret["metadata"].(map[string]interface{})["name"])
	}
}

func (h *ProjectHandler) Describe(params DescribeParams) {
	c := client.FromConfig()

	body, _, _ := c.Get("projects", params.Name)
	j, _ := yaml.JSONToYAML(body)

	fmt.Println(string(j))
}

func (h *ProjectHandler) Create(params CreateParams) {
	c := client.FromConfig()
	c.SetApiVersion(params.ApiVersion)

	body, _, _ := c.Post("projects", params.Resource)

	j, _ := yaml.JSONToYAML(body)

	fmt.Println(string(j))
}

func (h *ProjectHandler) Delete(params DeleteParams) {
	c := client.FromConfig()

	body, status, _ := c.Delete("projects", params.Name)

	if status == 200 {
		fmt.Printf("project \"%s\" deleted\n", params.Name)
	} else {
		fmt.Printf("failed to delete secret \"%s\"\n", params.Name)

		fmt.Println(string(body))
	}
}
