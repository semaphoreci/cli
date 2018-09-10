package cmd

import "github.com/ghodss/yaml"

func parse_yaml_to_map(data []byte) (map[string]interface{}, error) {
	m := make(map[string]interface{})

	err := yaml.Unmarshal(data, &m)

	return m, err
}
