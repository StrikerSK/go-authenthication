package server

import (
	"context"
	"fmt"
	userDomain "github.com/strikersk/user-auth/src/domain"
	userPorts "github.com/strikersk/user-auth/src/ports"
	"github.com/strikersk/user-auth/src/proto/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
)

type AuthorizationServer struct {
	service userPorts.IUserService
	auth.UnimplementedAuthorizationServiceServer
}

func NewAuthorizationServer(service userPorts.IUserService) *AuthorizationServer {
	return &AuthorizationServer{
		service: service,
	}
}

func (c *AuthorizationServer) RunServer() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Printf("Server init: %v\n", err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
	auth.RegisterAuthorizationServiceServer(grpcServer, c)

	if err = grpcServer.Serve(lis); err != nil {
		log.Printf("Server init: %v\n", err)
		os.Exit(1)
	}
}

func (c *AuthorizationServer) RegisterUser(ctx context.Context, in *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	user := userDomain.User{
		UserCredentials: userDomain.UserCredentials{
			Username: in.Username,
			Password: in.Password,
		},
	}

	fmt.Printf("Username: %s\nPasssword: %s\n", user.Username, user.Password)
	err := c.service.CreateUser(ctx, user)

	if err != nil {
		return &auth.RegisterResponse{
			Status: "Registration failed",
			Error:  err.Error(),
		}, err
	}

	response := &auth.RegisterResponse{
		Status: "User registered",
		Error:  "",
	}

	return response, nil
}

func (c *AuthorizationServer) LoginUser(ctx context.Context, in *auth.LoginRequest) (*auth.LoginResponse, error) {
	credentials := userDomain.UserCredentials{
		Username: in.Username,
		Password: in.Password,
	}

	token, err := c.service.LoginUser(ctx, credentials)

	if err != nil {
		return &auth.LoginResponse{
			Token: "",
			Error: err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	response := &auth.LoginResponse{
		Token: token,
		Error: "",
	}

	return response, nil
}
