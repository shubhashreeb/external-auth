package loginserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	_ "github.com/dgrijalva/jwt-go/v4"
	"github.com/shubhashreeb/external-auth/metrics"
	"github.com/shubhashreeb/external-auth/ratelimiter"
	"github.com/shubhashreeb/external-auth/redis"
	"gitlab.com/sitenet/svclib/logger"
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

type APIServer struct {
	cache       redis.Cache
	kCloak      *Keycloak
	ratelimiter *ratelimiter.RateLimiter
	metrics     metrics.Metrics
	logger      logger.Logger
}

func NewAPIServer() *APIServer {
	//Add to redis cache
	logger, err := logger.NewLogger()
	if err != nil {
		fmt.Printf("Could not instantiate log %s", err.Error())
	}
	logger.Info("Factory ...")
	config := redis.RedisConfig{Addrs: []string{getEnv("REDIS_ADDR", "redis:6379")}}
	cache, err := redis.NewRedisCache(&config)
	if err != nil {
		logger.Info("Error in connecting server")
	}
	kc := NewKeycloak(logger)
	m := *metrics.NewMetrics()

	return &APIServer{
		cache:       cache,
		kCloak:      kc,
		ratelimiter: ratelimiter.NewRateLimiter(logger),
		metrics:     m,
		logger:      logger,
	}
}

func (server *APIServer) login(w http.ResponseWriter, req *http.Request) {
	server.logger.Info("Login details is ", req.Body, "key: status", "okay", "statuscode", 200)

	//Increment prometheus counter for logoutReceived
	server.metrics.AddCounterStats("loginReceived", 1)

	var p LoginRequest
	err := json.NewDecoder(req.Body).Decode(&p)
	if err != nil {
		server.metrics.AddCounterStats("invalidLoginReceived", 1)
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger.Error(err.Error())
		return
	}

	server.logger.Info("Auth details is ", p)

	// check if the request is under rate limit
	opts := ratelimiter.RateLimitOpts{
		Domain:       "www.aaron.com",
		Path:         "/auth",
		Organization: "nshub",
		User:         "aaron",
	}
	if !server.ratelimiter.CheckIfRateUnderLimit(opts) {
		server.metrics.AddCounterStats("ratelimited", 1)
		server.logger.Info("Rate is over the set limit")
		http.Error(w, "rate exceeded", http.StatusBadRequest)
		return
	}

	res := server.kCloak.GetLoginToken(p)
	response, _ := json.Marshal(res)
	//server.logger.Info("Response :: ", response)

	loginRes := LoginResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		ExpiresIn:    res.ExpiresIn,
	}

	// Add the entry into cache - {Token}, {user-info}
	err = server.cache.Set(res.AccessToken, []byte("AccessToken"), 2000*time.Second)
	// err = server.cache.Set(loginRes.AccessToken, []byte("AccessToken"), 2000*time.Second)
	if err != nil {
		server.logger.Info("error setting key %s, Error : %v", loginRes.AccessToken, err)
	}

	server.metrics.AddCounterStats("cacheLoginToken", 1)

	server.logger.Info("Added token key to redis", loginRes.AccessToken, err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func (server *APIServer) logout(w http.ResponseWriter, req *http.Request) {
	server.logger.Info("Logout details is ", req.Body)

	//Increment prometheus counter for logoutReceived
	server.metrics.AddCounterStats("logoutReceived", 1)

	p := LogoutRequest{}
	// json.Marshal()
	err := json.NewDecoder(req.Body).Decode(&p)
	if err != nil {
		server.logger.Info("Error in json marshalling", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token := extractAuthToken(*req)
	server.logger.Info("Auth details is ", token)

	if token == "" {
		server.logger.Info("token is empty")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(""))
		return
	}

	res := server.kCloak.Logout(p)

	// Delete the token from the cache
	err = server.cache.Delete(token)

	if err != nil {
		server.logger.Info("Error in deleting token from cache")
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
		//server.logger.Info("error getting token from redis %s, Error : %v", token, err)
		server.logger.Info("error getting token from redis :", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	server.metrics.AddCounterStats("loginCacheHit", 1)
	server.logger.Info("Here is the token value fetched", bytes)

	// server.logger.Info("Added token key to redis", jwt.AccessToken, err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(""))
}

// fetch token from Keycloak
// create request with all the details such as clientId, secret etc to authenticate user
// once response is received, with token, store it in redis and return token to user
// this takes the http request and extract the authz header and
// extract the access token from the request and return it
// if not found it will return blank string
func extractAuthToken(req http.Request) string {
	authHeader := req.Header.Get("authorization")
	//server.logger.Info("Logout details is ", req.Body, authHeader)
	if authHeader == "" {
		fmt.Errorf("Authorization header missing")
		return ""
	}
	// Check if  Authorization header has Bearer token
	bearerToken := strings.Fields(authHeader)
	if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
		logger.Info("Invalid Auth header")
		return ""

	}
	logger.Info("Bearer token : ", bearerToken[1])
	if bearerToken[1] == "" {
		logger.Info("Token is missing")
	}
	return bearerToken[1]
}

func Serve() {
	server := NewAPIServer()
	go server.metrics.RunPrometheusServer()

	http.HandleFunc("/login", server.login)
	http.HandleFunc("/logout", server.logout)
	http.HandleFunc("/validate-token", server.validateToken)
	http.ListenAndServe(":8090", nil)
}
