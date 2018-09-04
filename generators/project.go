package generators

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/ghodss/yaml"
)

const project_template = `
apiVersion: v1alpha
kind: Project
metadata:
  name: %s
spec:
  repository:
    url: "%s"
`

func GenerateProjectYaml(name string, repo_url string) ([]byte, error) {
	return yaml.YAMLToJSON([]byte(fmt.Sprintf(project_template, name, repo_url)))
}

func GenerateProjectYamlFromRepoUrl(repo_url string) ([]byte, error) {
	name, err := ConstructProjectName(repo_url)

	if err != nil {
		return []byte{}, err
	}

	return GenerateProjectYaml(name, repo_url)
}

func ConstructProjectName(repo_url string) (string, error) {
	formats := []*regexp.Regexp{
		regexp.MustCompile(`git\@github\.com:.*\/(.*).git`),
		regexp.MustCompile(`git\@github\.com:.*\/(.*)`),
		regexp.MustCompile(`git\@github\.com:.*\/(.*)`),
	}

	for _, r := range formats {
		match := r.FindStringSubmatch(repo_url)

		if len(match) >= 2 {
			return match[1], nil
		}
	}

	errTemplate := "unknown git remote format '%s'.\n"
	errTemplate += "\n"
	errTemplate += "Format must be one of the following:\n"
	errTemplate += "  - git@github.com:/<owner>/<repo_name>.git\n"
	errTemplate += "  - git@github.com:/<owner>/<repo_name>\n"

	return "", errors.New(fmt.Sprintf(errTemplate, repo_url))
}
