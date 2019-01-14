package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func ContextList() ([]string, error) {
	res := []string{}

	m := make(map[string]interface{})

	err := viper.UnmarshalKey("contexts", &m)

	if err != nil {
		return []string{}, err
	}

	for k := range m {
		res = append(res, k)
	}

	return res, nil
}

func SetActiveContext(name string) {
	Set("active-context", name)
}

func GetActiveContext() string {
	if flag.Lookup("test.v") == nil {
		return viper.GetString("active-context")
	} else {
		return "org-semaphoretext-xyz"
	}
}

func GetAuth() string {
	if flag.Lookup("test.v") == nil {
		context := GetActiveContext()
		key_path := fmt.Sprintf("contexts.%s.auth.token", context)

		return Get(key_path)
	} else {
		return "123456789"
	}
}

func GetEditor() string {
	if flag.Lookup("test.v") == nil {
		editor := viper.GetString("editor")

		if editor != "" {
			return editor
		}

		editor = os.Getenv("EDITOR")

		if editor != "" {
			return editor
		}

		return "vim"
	} else {
		return "true" // Bash 'true' command, do nothing in tests
	}
}

func SetAuth(token string) {
	context := GetActiveContext()
	key_path := fmt.Sprintf("contexts.%s.auth.token", context)

	Set(key_path, token)
}

const unsetPublicSshKeyMsg = `Before creating a debug session job, configure the debug.PublicSshKey value.

Examples:

  # Configuring public ssh key with a literal
  sem config set debug.PublicSshKey "ssh-rsa AX3....DD"

  # Configuring public ssh key with a file
  sem config set debug.PublicSshKey "$(cat ~/.ssh/id_rsa.pub)"

  # Configuring public ssh key with your GitHub keys
  sem config set debug.PublicSshKey "$(curl -s https://github.com/<username>.keys)"
`

func GetPublicSshKeyForDebugSession() (string, error) {
	publicKey := viper.GetString("debug.PublicSshKey")

	if publicKey == "" {
		err := fmt.Errorf("Public SSH key for debugging is not configured.\n\n%s", unsetPublicSshKeyMsg)

		return "", err
	}

	return publicKey, nil
}

func GetHost() string {
	if flag.Lookup("test.v") == nil {
		context := GetActiveContext()
		key_path := fmt.Sprintf("contexts.%s.host", context)

		return Get(key_path)
	} else {
		return "org.semaphoretext.xyz"
	}
}

func SetHost(token string) {
	context := GetActiveContext()
	key_path := fmt.Sprintf("contexts.%s.host", context)

	Set(key_path, token)
}

func Set(key string, value string) {
	viper.Set(key, value)
	viper.WriteConfig()
}

func Get(key string) string {
	return viper.GetString(key)
}

func IsSet(key string) bool {
	return viper.IsSet(key)
}
