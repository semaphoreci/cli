package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/ghodss/yaml"
	"github.com/renderedtext/sem/client"
	"github.com/renderedtext/sem/cmd/utils"
	"github.com/renderedtext/sem/config"
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
	c := client.FromConfig()

	body, status, err := c.Get("secrets", params.Name)

	utils.Check(err)

	if status != 200 {
		utils.Fail(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	j, _ := yaml.JSONToYAML(body)

	content := fmt.Sprintf("# Editing Secrets/%s.\n# When you close the editor, the content will be updated on %s.\n\n%s", params.Name, config.GetHost(), j)

	tmpfile, err := ioutil.TempFile("", "sem-cli-edit-session")

	utils.Check(err)

	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		utils.Check(err)
	}

	if err := tmpfile.Close(); err != nil {
		utils.Check(err)
	}

	editor_path, err := exec.LookPath("vim")

	utils.Check(err)

	cmd := exec.Command(editor_path, tmpfile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Start()

	utils.Check(err)

	err = cmd.Wait()

	utils.Check(err)

	fmt.Printf("Secret '%s' updated.\n", params.Name)
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
