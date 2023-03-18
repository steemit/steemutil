package steemutil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type RequestData struct {
	Id      uint   `json:"id"`
	JsonRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
}

type ResponseData struct {
	Id      uint   `json:"id"`
	JsonRPC string `json:"jsonrpc"`
	Result  any    `json:"result"`
}

type IClient interface {
	Send(RequestData) (ResponseData, error)
}

type Client struct {
	Api     string
	Timeout uint
	Client  *http.Client
}

func (c *Client) Send(data RequestData) (res ResponseData, err error) {
	// Convert the data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	// Set up the HTTP request
	req, err := http.NewRequest("POST", c.Api, strings.NewReader(string(jsonData)))
	if err != nil {
		return
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")

	// TODO: timeout
	// Send the request and get the response
	resp, err := c.Client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Check http response status
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error:", resp.Status)
		return
	}

	// Parse response result
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(responseBody, &res)
	return
}

func GetClient(api string, timeout uint) (client IClient) {
	client = &Client{
		Api:     api,
		Timeout: timeout,
		Client:  &http.Client{},
	}
	return client
}
