package ratelimiter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"gitlab.com/sitenet/svclib/logger"
)

type RateLimiter struct {
	ipAddr string
	port   string
	client *http.Client
	url    string
	log    logger.Logger
}

type RateLimitReq struct {
	// The name of the rate limit IE: 'requests_per_second', 'gets_per_minute`
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Uniquely identifies this rate limit IE: 'ip:10.2.10.7' or 'account:123445'
	UniqueKey string `protobuf:"bytes,2,opt,name=unique_key,json=uniqueKey,proto3" json:"unique_key,omitempty"`
	// Rate limit requests optionally specify the number of hits a request adds to the matched limit. If Hit
	// is zero, the request returns the current limit, but does not increment the hit count.
	Hits int64 `protobuf:"varint,3,opt,name=hits,proto3" json:"hits,omitempty"`
	// The number of requests that can occur for the duration of the rate limit
	Limit int64 `protobuf:"varint,4,opt,name=limit,proto3" json:"limit,omitempty"`
	// The duration of the rate limit in milliseconds
	// Second = 1000 Milliseconds
	// Minute = 60000 Milliseconds
	// Hour = 3600000 Milliseconds
	Duration int64 `protobuf:"varint,5,opt,name=duration,proto3" json:"duration,omitempty"`

	Metadata map[string]string `protobuf:"bytes,9,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

type RateLimitResp struct {
	// The status of the rate limit.
	Status string `protobuf:"varint,1,opt,name=status,proto3,enum=pb.gubernator.Status" json:"status,omitempty"`
	// The currently configured request limit (Identical to RateLimitRequest.rate_limit_config.limit).
	Limit string `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
	// This is the number of requests remaining before the limit is hit.
	Remaining string `protobuf:"varint,3,opt,name=remaining,proto3" json:"remaining,omitempty"`
	// This is the time when the rate limit span will be reset, provided as a unix timestamp in milliseconds.
	ResetTime string `protobuf:"varint,4,opt,name=reset_time,json=resetTime,proto3" json:"reset_time,omitempty"`
	// Contains the error; If set all other values should be ignored
	Error string `protobuf:"bytes,5,opt,name=error,proto3" json:"error,omitempty"`
	// This is additional metadata that a client might find useful. (IE: Additional headers, corrdinator ownership, etc..)
	Metadata map[string]string `protobuf:"bytes,6,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

type GetRateLimitsReq struct {
	Requests []*RateLimitReq `protobuf:"bytes,1,rep,name=requests,proto3" json:"requests,omitempty"`
}

type GetRateLimitsResp struct {
	Responses []*RateLimitResp `protobuf:"bytes,1,rep,name=responses,proto3" json:"responses,omitempty"`
}

func NewRateLimiter(log logger.Logger) *RateLimiter {
	// shared HTTP transport and client for efficient connection reuse
	log.Info("Going to initialize the rate limiter")
	tr := &http.Transport{
		MaxIdleConns:          10,
		IdleConnTimeout:       15 * time.Second,
		ResponseHeaderTimeout: 15 * time.Second,
		DisableKeepAlives:     false,
	}
	httpClient := &http.Client{
		Transport: tr,
		// do not follow redirects
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	//Gubernator address
	ipAddr := "192.168.86.211"
	port := "30080" //"9080"
	return &RateLimiter{
		ipAddr: ipAddr,
		port:   port,
		client: httpClient,
		url:    fmt.Sprintf("http://%s:%s/v1/GetRateLimits", ipAddr, port),
		log:    log,
	}
}

type RateLimitOpts struct {
	Domain       string
	Path         string
	Organization string
	User         string
}

func (e RateLimiter) CheckIfRateUnderLimit(opts RateLimitOpts) bool {
	e.log.Info("Going to check with Gubernator ...")

	getReq := GetRateLimitsReq{
		Requests: make([]*RateLimitReq, 0),
	}

	req1 := []*RateLimitReq{
		// rate limit for global
		{
			Name:      "requests_per_sec",
			UniqueKey: fmt.Sprintf(opts.Domain),
			Hits:      1,
			Limit:     100,
			Duration:  5000,
		},
		// limit for path
		{
			Name:      "requests_per_sec",
			UniqueKey: fmt.Sprintf("%s-%s", opts.Domain, opts.Path),
			Hits:      1,
			Limit:     10,
			Duration:  5000,
		},
		// user limit
		{
			Name:      "requests_per_sec",
			UniqueKey: opts.User,
			Hits:      1,
			Limit:     5,
			Duration:  5000,
		},
		// user path limit
		{
			Name:      "requests_per_sec",
			UniqueKey: fmt.Sprintf("%s-%s-%s", opts.Domain, opts.Path, opts.User),
			Hits:      1,
			Limit:     1,
			Duration:  5000,
		},
	}

	getReq.Requests = req1 //append(getReq.Requests, &req1)
	jsonData, err := json.Marshal(getReq)
	requestPayload := []byte(jsonData)
	req, err := http.NewRequest(
		http.MethodPost,
		e.url,
		bytes.NewBuffer(requestPayload))

	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		os.Exit(1)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		e.log.Info("Fatiled to fetch body : ", err)
	}

	e.log.Info("Response body : ", string(body))
	var rateLimitResp *GetRateLimitsResp
	err = json.Unmarshal(body, &rateLimitResp)
	e.log.Info("Response respnse : ", rateLimitResp, err)

	// loop over all the responses from the array
	// if any of the response is not underlimit, return false
	for cnt, resp := range rateLimitResp.Responses {
		if resp.Status != "UNDER_LIMIT" {
			e.log.Info("-- RATE LIMIT HIT -- for ", cnt)
			return false
		}
	}

	e.log.Info("Remaining limit :: ", rateLimitResp.Responses[0].Remaining)
	return true
}
