package loginserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	_ "github.com/dgrijalva/jwt-go/v4"
	"github.com/shubhashreeb/external-auth/keycloak"
	"github.com/shubhashreeb/external-auth/redis"
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

// it receives the logout Request
// 1st -> validate if this contains the token in its body
// Get the Authorization header from the request
// 2nd -> sends logout request to keycloak
// 3rd -> Remove the redis cache
// return the response
func (server APIServer) logout(w http.ResponseWriter, req *http.Request) {
	authHeader := req.Header.Get("authorization")
	//fmt.Println("Logout details is ", req.Body, authHeader)
	if authHeader == "" {
		fmt.Errorf("Authorization header missing")
	}
	// Check if  Authorization header has Bearer token

	bearerToken := strings.Fields(authHeader)
	if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
		fmt.Errorf("Invalid Auth header")
	}
	fmt.Println("Bearer token : ", bearerToken[1])

}

func (server APIServer) login(w http.ResponseWriter, req *http.Request) {
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

	jwt, err := keycloak.
		NewMyHTTPClient().
		SendRequest(*keycloak.GetNewKeycloakConfig(), login)

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

	err = server.cache.Set(jwt.AccessToken, []byte("AccessToken"), 2000*time.Second)
	if err != nil {
		fmt.Println("error setting key %s, Error : %v", jwt.AccessToken, err)
	}

	fmt.Println("Added token key to redis", jwt.AccessToken, err)
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

// validate if token is present in request
// if token is present, it will check with redis if the token is valid
// if token is valid, it will return the response

func (server APIServer) validateToken(w http.ResponseWriter, req *http.Request) {
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
	bytes, err := server.cache.Get(bearerToken[1])
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

func Server() {
	apiServer := NewAPIServer()
	// keycloakClient = NewKeycloak().Setup()
	http.HandleFunc("/login", apiServer.login)
	http.HandleFunc("/logout", apiServer.logout)
	http.HandleFunc("/validate-token", apiServer.validateToken)

	fmt.Println("Running the login server on port: 8090")
	http.ListenAndServe(":8090", nil)
}
