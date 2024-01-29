package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Nerzal/gocloak/v13"
	_ "github.com/dgrijalva/jwt-go/v4"
	"github.com/shubhashreeb/external-auth/redis"
)

var keycloakClient *Keycloak

type JWT struct {
	AccessToken      string `json:"access_token"`
	IDToken          string `json:"id_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expiresIn"`
}

type Keycloak struct {
	keycloak     *gocloak.GoCloak
	address      string
	clientId     string // clientId specified in Keycloak
	clientSecret string // client secret specified in Keycloak
	realm        string // realm specified in Keycloak
}

func NewKeycloak() *Keycloak {
	return &Keycloak{
		address:      "http://192.168.86.211:32088/",
		clientId:     "auth-svc",
		clientSecret: "PIvkN94ImvghqZmv0vJUO2PElHtWYXsY", //"eQxQdPudNTyi8rv5L3Tgs1SO5byD6vNB",
		realm:        "nshub",
	}
}

type APIServer struct {
	cache redis.Cache
}

func NewAPIServer() *APIServer {
	//Add to redis cache
	config := redis.RedisConfig{Addrs: []string{"192.168.86.211:32379"}}
	cache, err := redis.NewRedisCache(&config)
	if err != nil {
		fmt.Println("Error in connecting server")
	}
	return &APIServer{
		cache: cache,
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
	res := keycloakClient.GetLoginToken(p)
	response, _ := json.Marshal(res)
	fmt.Println("Response :: ", response)
	// Add the entry into cache - {Token}, {user-info}

	//fmt.Fprintf(w, "hello\n")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func logout(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Logout details is ", req.Body)
	var p LogoutRequest
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
	res := keycloakClient.Logout(p)
	// response, _ := json.Marshal(res)
	// fmt.Println("Response :: ", response)
	// Add the entry into cache - {Token}, {user-info}

	//fmt.Fprintf(w, "hello\n")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(res))
}

func validateToken(w http.ResponseWriter, req *http.Request) {
	authHeader := req.Header.Get("authorization")
	//fmt.Println("Logout details is ", req.Body, authHeader)
	if authHeader == "" {
		fmt.Errorf("Authorization header missing")
	}
	// Check if  Authorization header has Bearer token

	bearerToken := strings.Fields(authHeader)
	if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
		fmt.Println("Invalid Auth header")
	}
	fmt.Println("Bearer token : ", bearerToken[1])

	if bearerToken[1] == "" {
		fmt.Println("Token is missing")

	}
	bytes, err := NewAPIServer().cache.Get(bearerToken[1])
	if err != nil {
		fmt.Println("error getting token from redis %s, Error : %v", bearerToken[1], err)
		w.WriteHeader(http.StatusUnauthorized)
	}

	fmt.Println("Here is the token value fetched", string(bytes))

	// fmt.Println("Added token key to redis", jwt.AccessToken, err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(""))
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

func (k *Keycloak) Logout(req LogoutRequest) string {
	//jwt, err := k.gocloak.Login(context.Background(),
	err := k.keycloak.Logout(context.Background(),
		k.clientId,
		k.clientSecret,
		k.realm,
		req.RefreshToken,
	)

	fmt.Println("Logout response received", err)

	if err != nil {
		// http.Error(w, err.Error(), http.StatusForbidden)
		return fmt.Sprintf("invalid with error", err)
	}
	// fmt.Println("Here is the token response", jwt)

	// TODO - Add to cache

	return "Sucessfully logged out"
}

func main() {
	keycloakClient = NewKeycloak().Setup()
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/validate-token", validateToken)
	http.ListenAndServe(":8090", nil)
}
