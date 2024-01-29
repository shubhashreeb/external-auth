package loginserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

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

type KeycloakConfig struct {
	address      string
	clientId     string // clientId specified in Keycloak
	clientSecret string // client secret specified in Keycloak
	realm        string // realm specified in Keycloak
}

type Keycloak struct {
	client *gocloak.GoCloak
	config KeycloakConfig
}

func NewKeycloak() *Keycloak {
	c := KeycloakConfig{
		address:      "http://192.168.86.211:32088/",
		clientId:     "auth-svc",
		clientSecret: "PIvkN94ImvghqZmv0vJUO2PElHtWYXsY", //"eQxQdPudNTyi8rv5L3Tgs1SO5byD6vNB",
		realm:        "nshub",
	}
	return &Keycloak{
		config: c,
		client: gocloak.NewClient(c.address),
	}
}

type APIServer struct {
	cache  redis.Cache
	kCloak *Keycloak
}

func NewAPIServer() *APIServer {
	//Add to redis cache
	config := redis.RedisConfig{Addrs: []string{"192.168.86.211:32379"}}
	cache, err := redis.NewRedisCache(&config)
	if err != nil {
		fmt.Println("Error in connecting server")
	}
	kc := NewKeycloak()
	return &APIServer{
		cache:  cache,
		kCloak: kc,
	}
}

func (server *APIServer) login(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Login details is ", req.Body)
	var p LoginRequest
	err := json.NewDecoder(req.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("Auth details is ", p)

	// for name, headers := range req.Header {
	// 	for _, h := range headers {
	// 		fmt.Println("%v: %v\n", name, h)
	// 	}
	// }

	res := server.kCloak.GetLoginToken(p)
	response, _ := json.Marshal(res)
	fmt.Println("Response :: ", response)

	loginRes := LoginResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		ExpiresIn:    res.ExpiresIn,
	}

	// Add the entry into cache - {Token}, {user-info}
	err = server.cache.Set(res.AccessToken, []byte("AccessToken"), 2000*time.Second)
	// err = server.cache.Set(loginRes.AccessToken, []byte("AccessToken"), 2000*time.Second)
	if err != nil {
		fmt.Println("error setting key %s, Error : %v", loginRes.AccessToken, err)
	}

	fmt.Println("Added token key to redis", loginRes.AccessToken, err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func (server *APIServer) logout(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Logout details is ", req.Body)
	p := LogoutRequest{}
	// json.Marshal()
	err := json.NewDecoder(req.Body).Decode(&p)
	if err != nil {
		fmt.Println("Error in json marshalling", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token := extractAuthToken(*req)
	fmt.Println("Auth details is ", token)

	if token == "" {
		fmt.Println("token is empty")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(""))
		return
	}

	res := server.kCloak.Logout(p)

	// Delete the token from the cache
	err = server.cache.Delete(token)

	if err != nil {
		fmt.Println("Error in deleting token from cache")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(res))
}

func (server *APIServer) validateToken(w http.ResponseWriter, req *http.Request) {
	token := extractAuthToken(*req)
	//check if the token exists in Redis
	bytes, err := server.cache.Get(token)
	if err != nil {
		//fmt.Println("error getting token from redis %s, Error : %v", token, err)
		fmt.Println("error getting token from redis :", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	fmt.Println("Here is the token value fetched", bytes)

	// fmt.Println("Added token key to redis", jwt.AccessToken, err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(""))
}

// fetch token from Keycloak
// create request with all the details such as clientId, secret etc to authenticate user
// once response is received, with token, store it in redis and return token to user

func (k *Keycloak) GetLoginToken(req LoginRequest) *LoginResponse {
	jwt, err := k.client.Login(context.Background(),
		k.config.clientId,
		k.config.clientSecret,
		k.config.realm,
		req.Username,
		req.Password,
	)

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
	fmt.Println("Logout req received with refresh token", req.RefreshToken)

	err := k.client.Logout(context.Background(),
		k.config.clientId,
		k.config.clientSecret,
		k.config.realm,
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

// this takes the http request and extract the authz header and
// extract the access token from the request and return it
// if not found it will return blank string
func extractAuthToken(req http.Request) string {
	authHeader := req.Header.Get("authorization")
	//fmt.Println("Logout details is ", req.Body, authHeader)
	if authHeader == "" {
		fmt.Errorf("Authorization header missing")
		return ""
	}
	// Check if  Authorization header has Bearer token
	bearerToken := strings.Fields(authHeader)
	if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
		fmt.Println("Invalid Auth header")
		return ""

	}
	fmt.Println("Bearer token : ", bearerToken[1])
	if bearerToken[1] == "" {
		fmt.Println("Token is missing")
	}
	return bearerToken[1]
}

func Serve() {
	server := NewAPIServer()
	http.HandleFunc("/login", server.login)
	http.HandleFunc("/logout", server.logout)
	http.HandleFunc("/validate-token", server.validateToken)
	http.ListenAndServe(":8090", nil)
}
