package handler

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ghodss/yaml"
	"github.com/renderedtext/sem/client"
	"github.com/renderedtext/sem/cmd/utils"
)

type SecretHandler struct {
}

func (h *SecretHandler) Get(params GetParams) {
	secretList, err := client.ListSecrets()

	utils.Check(err)

	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)

	fmt.Fprintln(w, "NAME\tAGE")

	for _, secret := range secretList.Secrets {
		fmt.Fprintf(w, "%s\t%s\n", secret.Metadata.Name, RelativeAgeForHumans(secret.Metadata.UpdateTime))
	}

	w.Flush()
}

func (h *SecretHandler) Describe(params DescribeParams) {
	c := client.FromConfig()
	c.SetApiVersion("v1beta")

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
	c.SetApiVersion("v1beta")

	s, err := client.InitSecretFromYaml(params.Resource)

	utils.Check(err)

	err = s.Create()

	utils.Check(err)

	fmt.Printf("Secret '%s' updated.\n", s.Metadata.Name)
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
	err := client.DeleteSecret(params.Name)

	utils.Check(err)

	fmt.Printf("secret \"%s\" deleted\n", params.Name)
}
