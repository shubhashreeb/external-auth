package client

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpClient(t *testing.T) {
	c := NewMyHTTPClient()
	url := "192.168.86.211"
	port := "32088"
	protocol := "openid-connect"
	realm := "nshub"
	client_id := "auth-svc"
	client_secret := "3nUB5EUG3fenYFa9xqnF376PLFZHWxFV"
	username := "aaron"
	password := "aaron"
	res, err := c.SendRequest(url, port, protocol, realm, client_id, client_secret, username, password)
	assert.Nil(t, err)
	fmt.Println("Res body: ", res)
}
