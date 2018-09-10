package handler

import (
	"encoding/json"
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/renderedtext/sem/client"
	"github.com/renderedtext/sem/cmd/utils"
)

type SecretHandler struct {
}

func (h *SecretHandler) Get(params GetParams) {
	c := client.FromConfig()

	body, status, err := c.List("secrets")

	utils.Check(err)

	if status != 200 {
		utils.Fail(fmt.Sprintf("http status %d received from upstream", status))
	}

	var secrets []map[string]interface{}

	json.Unmarshal([]byte(body), &secrets)

	fmt.Println("NAME")

	for _, secret := range secrets {
		fmt.Println(secret["metadata"].(map[string]interface{})["name"])
	}
}

func (h *SecretHandler) Describe(params DescribeParams) {
	c := client.FromConfig()

	body, status, err := c.Get("secrets", params.Name)

	utils.Check(err)

	if status != 200 {
		utils.Fail(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	j, _ := yaml.JSONToYAML(body)

	fmt.Println(string(j))
}

func (h *SecretHandler) Create(params CreateParams) {
	c := client.FromConfig()
	c.SetApiVersion(params.ApiVersion)

	body, status, err := c.Post("secrets", params.Resource)

	utils.Check(err)

	if status != 200 {
		utils.Fail(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	j, _ := yaml.JSONToYAML(body)

	fmt.Println(string(j))
}

func (h *SecretHandler) Edit(params EditParams) {
	secret, err := client.GetSecret(params.Name)

	utils.Check(err)

	content, err := secret.ToYaml()

	utils.Check(err)

	new_content, err := EditYamlInEditor(secret.ObjectName(), string(content))

	utils.Check(err)

	updated_secret, err := client.InitSecretFromYaml([]byte(new_content))

	utils.Check(err)

	err = updated_secret.Update()

	utils.Check(err)

	fmt.Printf("Secret '%s' updated.\n", secret.Metadata.Name)
}

func (h *SecretHandler) Delete(params DeleteParams) {
	c := client.FromConfig()

	body, status, err := c.Delete("secrets", params.Name)

	utils.Check(err)

	if status != 200 {
		utils.Fail(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	fmt.Printf("secret \"%s\" deleted\n", params.Name)
}
