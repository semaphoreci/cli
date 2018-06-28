package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/spf13/viper"
)

type Client struct {
	authToken  string
	host       string
	apiVersion string
}

func FromConfig() Client {
	authToken := viper.GetString("auth.token")
	host := viper.GetString("host")
	apiVersion := "v1alpha"

	if authToken == "" {
		fmt.Println("Auth Token is not configured.")
		fmt.Println("Run the following command to set it:")
		fmt.Println("")
		fmt.Println("  sem config set auth.token [AUTH_TOKEN]")
		fmt.Println("")

		os.Exit(1)
	}

	if host == "" {
		fmt.Println("Api host is not configured.")
		fmt.Println("Run the following command to set it:")
		fmt.Println("")
		fmt.Println("  sem config set host [API_HOST]")
		fmt.Println("")

		os.Exit(1)
	}

	return New(authToken, host, apiVersion)
}

func New(authToken string, host string, apiVersion string) Client {
	return Client{authToken, host, apiVersion}
}

func (c *Client) SetApiVersion(apiVersion string) *Client {
	c.apiVersion = apiVersion

	return c
}

func (c *Client) Get(kind string, name string) ([]byte, int, error) {
	url := fmt.Sprintf("https://%s/api/%s/%s/%s", c.host, c.apiVersion, kind, name)

	log.Println(url)

	req, err := http.NewRequest("GET", url, nil)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", c.authToken))

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return []byte(""), 0, err
	}

	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)

	body, err := ioutil.ReadAll(resp.Body)

	log.Println(string(body))

	return body, resp.StatusCode, err
}

func (c *Client) List(kind string) ([]byte, int, error) {
	url := fmt.Sprintf("https://%s/api/%s/%s", c.host, c.apiVersion, kind)

	log.Println(url)

	req, err := http.NewRequest("GET", url, nil)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", c.authToken))

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return []byte(""), 0, err
	}

	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)

	body, err := ioutil.ReadAll(resp.Body)

	log.Println(string(body))

	return body, resp.StatusCode, err
}

func (c *Client) Delete(kind string, name string) ([]byte, int, error) {
	url := fmt.Sprintf("https://%s/api/%s/%s/%s", c.host, c.apiVersion, kind, name)

	log.Println(url)

	req, err := http.NewRequest("DELETE", url, nil)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", c.authToken))

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return []byte(""), 0, err
	}

	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)

	body, err := ioutil.ReadAll(resp.Body)

	log.Println(string(body))

	return body, resp.StatusCode, err
}

func (c *Client) Post(kind string, resource []byte) ([]byte, int, error) {
	url := fmt.Sprintf("https://%s/api/%s/%s", c.host, c.apiVersion, kind)

	log.Println(url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(resource))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", c.authToken))

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return []byte(""), 0, err
	}

	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)

	body, err := ioutil.ReadAll(resp.Body)

	log.Println(string(body))

	return body, resp.StatusCode, err
}
