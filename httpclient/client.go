package client

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"time"
)

//create a Http client, create a request by constructing the url, body, make the request using clients Do method and read the response body

type HTTPClient struct {
	client *http.Client
}

// NewHTTPClient creates a new instance of HTTPClient
func NewMyHTTPClient() *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SendRequest creates a request by constructing the url, body, make the request using clients Do method and read the response body
func (c *HTTPClient) SendRequest(url string, port string, protocol string, realm string, client_id string, client_secret string, username string, password string) (*http.Response, error) {

	requestURL := fmt.Sprintf("http://%s:%s/realms/%s/protocol/%s/token", url, port, realm, protocol)

	//requestURL := "http://192.168.86.211:32088/realms/nshub/protocol/openid-connect/token"
	jsonBody := []byte(fmt.Sprintf(
		"grant_type=password&client_id=%s&client_secret=%s&username=%s&password=%s",
		client_id,
		client_secret,
		username,
		password))

	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := c.client.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		os.Exit(1)
		return nil, err
	}
	fmt.Println("The response body :", res)
	return res, nil
}
