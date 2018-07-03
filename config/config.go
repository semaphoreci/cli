package config

import (
	"fmt"

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
	return viper.GetString("active-context")
}

func GetAuth() string {
	context := GetActiveContext()
	key_path := fmt.Sprintf("contexts.%s.auth.token", context)

	return Get(key_path)
}

func SetAuth(token string) {
	context := GetActiveContext()
	key_path := fmt.Sprintf("contexts.%s.auth.token", context)

	Set(key_path, token)
}

func GetHost() string {
	context := GetActiveContext()
	key_path := fmt.Sprintf("contexts.%s.host", context)

	return Get(key_path)
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
