package handler

import (
	"fmt"
	"os"
	"text/tabwriter"

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
	s, err := client.GetSecret(params.Name)

	utils.Check(err)

	body, err := s.ToYaml()

	fmt.Print(string(body))
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

func (h *SecretHandler) Apply(params ApplyParams) {
	c := client.FromConfig()
	c.SetApiVersion("v1beta")

	s, err := client.InitSecretFromYaml(params.Resource)

	utils.Check(err)

	err = s.Update()

	utils.Check(err)

	fmt.Printf("Secret '%s' updated.\n", s.Metadata.Name)
}

func (h *SecretHandler) Delete(params DeleteParams) {
	err := client.DeleteSecret(params.Name)

	utils.Check(err)

	fmt.Printf("secret \"%s\" deleted\n", params.Name)
}
