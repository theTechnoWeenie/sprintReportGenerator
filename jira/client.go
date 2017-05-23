package jira

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Credentials struct {
	User     string
	Password string
}

func (credentials *Credentials) Init(user string, password string) {
	credentials.User = strings.Replace(strings.Replace(user, "\n", "", -1), "\r", "", -1)
	credentials.Password = strings.Replace(password, "\\!", "!", -1)
}

func (credentials *Credentials) GetEncoded() string {
	encoded := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", credentials.User, credentials.Password)))
	return encoded
}

type Client struct {
	BaseURL     string
	credentials Credentials
	client      *http.Client
}

func NewClient(baseURL string, creds Credentials) Client {
	return Client{baseURL, creds, &http.Client{}}
}

func (client *Client) SetCredentials(credentials Credentials) {
	client.credentials = credentials
}

func (client *Client) Get(path string, query map[string]string) []byte {
	path_and_query := path + "?"
	for param, value := range query {
		path_and_query += fmt.Sprintf("%s=%s&", param, value)
	}
	resource := fmt.Sprintf("%s%s", client.BaseURL, path_and_query)
	req, _ := http.NewRequest("GET", resource, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", client.credentials.GetEncoded()))
	resp, err := client.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return body
}
