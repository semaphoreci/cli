package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/semaphoreci/cli/config"
)

type BaseClient struct {
	authToken  string
	host       string
	apiVersion string
}

func NewBaseClientFromConfig() BaseClient {
	host := config.GetHost()
	authToken := config.GetAuth()
	apiVersion := "v1alpha"

	if authToken == "" || host == "" {
		fmt.Println("Connection to Semaphore is not established.")
		fmt.Println("Run the following command to connect to Semaphore:")
		fmt.Println("")
		fmt.Println("  sem connect [HOST] [TOKEN]")
		fmt.Println("")

		os.Exit(1)
	}

	return NewBaseClient(authToken, host, apiVersion)
}

func NewBaseClient(authToken string, host string, apiVersion string) BaseClient {
	return BaseClient{authToken, host, apiVersion}
}

func (c *BaseClient) SetApiVersion(apiVersion string) *BaseClient {
	c.apiVersion = apiVersion

	return c
}

func (c *BaseClient) Get(kind string, name string) ([]byte, int, error) {
	url := fmt.Sprintf("https://%s/api/%s/%s/%s", c.host, c.apiVersion, kind, name)

	log.Printf("GET %s\n", url)

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

func (c *BaseClient) List(kind string) ([]byte, int, error) {
	url := fmt.Sprintf("https://%s/api/%s/%s", c.host, c.apiVersion, kind)

	log.Printf("GET %s\n", url)

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

func (c *BaseClient) ListWithParams(kind string, query url.Values) ([]byte, int, error) {
	url := fmt.Sprintf("https://%s/api/%s/%s?%s", c.host, c.apiVersion, kind, query.Encode())

	log.Printf("GET %s\n", url)

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

func (c *BaseClient) Delete(kind string, name string) ([]byte, int, error) {
	url := fmt.Sprintf("https://%s/api/%s/%s/%s", c.host, c.apiVersion, kind, name)

	log.Printf("DELETE %s\n", url)

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

func (c *BaseClient) PostAction(kind, item, action string, resource []byte) ([]byte, int, error) {
	kindItemAction := fmt.Sprintf("%s/%s/%s", kind, item, action)
	return c.Post(kindItemAction, resource)
}

func (c *BaseClient) Post(kind string, resource []byte) ([]byte, int, error) {
	return c.PostHeaders(kind, resource, make(map[string]string))
}

func (c *BaseClient) PostHeaders(kind string, resource []byte, headers map[string]string) ([]byte, int, error) {
	url := fmt.Sprintf("https://%s/api/%s/%s", c.host, c.apiVersion, kind)

	log.Printf("POST %s\n", url)
	log.Println("Resource", string(resource))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(resource))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", c.authToken))

	for k, v := range headers {
		req.Header.Set(k, v)
	}

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

func (c *BaseClient) Patch(kind string, name string, resource []byte) ([]byte, int, error) {
	url := fmt.Sprintf("https://%s/api/%s/%s/%s", c.host, c.apiVersion, kind, name)

	log.Printf("PATCH %s\n", url)

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(resource))

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
