package keycloak

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
func (c *HTTPClient) SendRequest(kConf KeycloakConfig, login LoginRequest) (*JWT, error) {
	requestURL := fmt.Sprintf("http://%s:%s/realms/%s/protocol/%s/token", kConf.Address, kConf.Port, kConf.Realm, kConf.Protocol)
	jsonBody := []byte(fmt.Sprintf(
		"grant_type=password&client_id=%s&client_secret=%s&username=%s&password=%s",
		kConf.ClientId, kConf.ClientSecret, login.Username, login.Password))

	bodyReader := bytes.NewReader(jsonBody)
	req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := c.client.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		return nil, err
	}

	if res.StatusCode != 200 {
		fmt.Printf("http response code is not OK: %s\n", res.Status)
		return nil, fmt.Errorf("http response code is not OK: %s\n", res.Status)
	}

	bytes, e := io.ReadAll(res.Body)
	if e != nil {
		fmt.Printf("Error in reading the response body: %s\n", e)
		return nil, e
	}
	fmt.Println("The response body with JWT token :", string(bytes))

	// unmarshal the respnse and extract JWT from response body
	jwt := &JWT{}
	e = json.Unmarshal(bytes, jwt)
	if e != nil {
		fmt.Printf("Error in unmarshalling the token : %s\n", e)
		return nil, e
	}

	return jwt, nil

}
