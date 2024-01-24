package main

import (
	"fmt"
	"log"
	"net"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/gogo/googleapis/google/rpc"
	"golang.org/x/net/context"
	rpcstatus "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
)

type AuthorizationServer struct{}

func (a *AuthorizationServer) Check(ctx context.Context, req *auth.CheckRequest) (*auth.CheckResponse, error) {
	log.Println(">>> Server Performing authorization check!")
	method := req.Attributes.Request.Http.Method
	path := req.Attributes.Request.Http.Path
	authHeader, ok := req.Attributes.Request.Http.Headers["authorization"]

	if !ok {
		fmt.Println("failed to receive headers")
	}

	fmt.Println("Here are the request headers", req)
	fmt.Println("Here are the request info", method, path, authHeader)

	return &auth.CheckResponse{
		Status: &rpcstatus.Status{
			Code: int32(rpc.OK),
		},
		HttpResponse: &auth.CheckResponse_OkResponse{
			OkResponse: &auth.OkHttpResponse{
				Headers: []*core.HeaderValueOption{
					{
						Header: &core.HeaderValue{
							Key:   "x-prism-authorized",
							Value: "allowed",
						},
					},
				},
			},
		},
	}, nil

	/*
		return &auth.CheckResponse{
					Status: &rpcstatus.Status{
						Code: int32(rpc.PERMISSION_DENIED), // Code: int32(rpc.UNAUTHENTICATED),
					},
					HttpResponse: &auth.CheckResponse_DeniedResponse{
						DeniedResponse: &auth.DeniedHttpResponse{
							Status: &envoy_type.HttpStatus{
								Code: envoy_type.StatusCode_Unauthorized,
							},
							Body: "PERMISSION_DENIED",
						},
					},
				}, nil
	*/

}

func main() {
	fmt.Println("Running the auth server ...")
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	auth.RegisterAuthorizationServer(s, &AuthorizationServer{})

	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
