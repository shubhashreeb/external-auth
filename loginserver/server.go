package loginserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/dgrijalva/jwt-go/v4"
	"github.com/shubhashreeb/external-auth/keycloak"
)

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int    `json:"expiresIn"`
}

// func NewKeycloak() *Keycloak {
// 	return &Keycloak{
// 		address:      "http://192.168.86.211:32088/",
// 		clientId:     "auth-svc",
// 		clientSecret: "3nUB5EUG3fenYFa9xqnF376PLFZHWxFV", //"eQxQdPudNTyi8rv5L3Tgs1SO5byD6vNB",
// 		realm:        "nshub",
// 	}
// }

// func (k *Keycloak) Setup() *Keycloak {
// 	k.keycloak = gocloak.NewClient(k.address)
// 	return k
// }

func login(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Login details is ", req.Body)
	login := keycloak.LoginRequest{}
	err := json.NewDecoder(req.Body).Decode(&login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("Auth details is ", login)

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}

	jwt, err := keycloak.NewMyHTTPClient().SendRequest(*keycloak.GetNewKeycloakConfig(), login)

	loginRes := LoginResponse{
		AccessToken:  jwt.AccessToken,
		RefreshToken: jwt.RefreshToken,
		ExpiresIn:    jwt.ExpiresIn,
	}
	loginRes.ExpiresIn = jwt.ExpiresIn
	json, err := json.Marshal(loginRes)
	if err != nil {
		fmt.Println("Failed to Marshal the login request", err)
		_, _ = w.Write([]byte("internal error"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(json)
}

// fetch token from Keycloak
// create request with all the details such as clientId, secret etc to authenticate user
// once response is received, with token, store it in redis and return token to user

// func (k *Keycloak) GetLoginToken(req LoginRequest) *LoginResponse {
// 	//jwt, err := k.gocloak.Login(context.Background(),
// 	jwt, err := k.keycloak.Login(context.Background(),
// 		k.clientId,
// 		k.clientSecret,
// 		k.realm,
// 		req.Username,
// 		req.Password)

// 	fmt.Println("Login request received", jwt, err)

// 	if err != nil {
// 		// http.Error(w, err.Error(), http.StatusForbidden)
// 		return &LoginResponse{}
// 	}
// 	fmt.Println("Here is the token response", jwt)

// 	// TODO - Add to cache

// 	return &LoginResponse{
// 		AccessToken:  jwt.AccessToken,
// 		RefreshToken: jwt.RefreshToken,
// 		ExpiresIn:    jwt.ExpiresIn,
// 	}
// }

func Server() {
	// keycloakClinet = NewKeycloak().Setup()
	http.HandleFunc("/login", login)
	fmt.Println("Running the login server on port: 8090")
	http.ListenAndServe(":8090", nil)
}
