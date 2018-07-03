package config

import (
	"fmt"
	"github.com/spf13/viper"
)

func ContextList() ([]string, error) {
	m := make(map[string]interface{})

	err := viper.UnmarshalKey("contexts", &m)

	if err != nil {
		return []string{}, err
	}

   for k := range m {
	   if k == GetActiveContext() {
		   fmt.Print("* ")
		   fmt.Println(k)
	   } else {
		   fmt.Print("  ")
		   fmt.Println(k)
	   }
    }

	return []string{}, nil
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
