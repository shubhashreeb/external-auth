package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Nerzal/gocloak/v7"
	_ "github.com/dgrijalva/jwt-go/v4"
)

var keycloakClinet *Keycloak

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int    `json:"expiresIn"`
}

func NewKeycloak() *Keycloak {
	return &Keycloak{
		address:      "http://192.168.86.211:32088/",
		clientId:     "auth-svc",
		clientSecret: "3nUB5EUG3fenYFa9xqnF376PLFZHWxFV", //"eQxQdPudNTyi8rv5L3Tgs1SO5byD6vNB",
		realm:        "nshub",
	}
}

func (k *Keycloak) Setup() *Keycloak {
	k.keycloak = gocloak.NewClient(k.address)
	return k
}

func login(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Login details is ", req.Body)
	var p LoginRequest
	err := json.NewDecoder(req.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("Auth details is ", p)

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
	res := keycloakClinet.GetLoginToken(p)
	response, _ := json.Marshal(res)
	fmt.Println("Response  received :: ", response)
	// Add the entry into cache - {Token}, {user-info}

	//fmt.Fprintf(w, "hello\n")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

// fetch token from Keycloak
// create request with all the details such as clientId, secret etc to authenticate user
// once response is received, with token, store it in redis and return token to user

func (k *Keycloak) GetLoginToken(req LoginRequest) *LoginResponse {
	//jwt, err := k.gocloak.Login(context.Background(),
	jwt, err := k.keycloak.Login(context.Background(),
		k.clientId,
		k.clientSecret,
		k.realm,
		req.Username,
		req.Password)

	fmt.Println("Login request received", jwt, err)

	if err != nil {
		// http.Error(w, err.Error(), http.StatusForbidden)
		return &LoginResponse{}
	}
	fmt.Println("Here is the token response", jwt)

	// TODO - Add to cache

	return &LoginResponse{
		AccessToken:  jwt.AccessToken,
		RefreshToken: jwt.RefreshToken,
		ExpiresIn:    jwt.ExpiresIn,
	}
}

func main() {
	keycloakClinet = NewKeycloak().Setup()
	http.HandleFunc("/login", login)
	http.ListenAndServe(":8090", nil)
}
