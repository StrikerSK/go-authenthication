package server

import (
	"context"
	"github.com/strikersk/user-auth/src/domain"
	userPorts "github.com/strikersk/user-auth/src/ports"
	"github.com/strikersk/user-auth/src/proto/auth"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

type AuthorizationServer struct {
	service userPorts.IUserService
	auth.UnimplementedAuthorizationServiceServer
}

func (c *AuthorizationServer) RegisterUser(ctx context.Context, in *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	user := domain.User{
		UserCredentials: domain.UserCredentials{
			Username: in.Username,
			Password: in.Password,
		},
	}

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

func CreateAuthorizationServer(service userPorts.IUserService) {
	lis, err := net.Listen("tcp", "9000")
	if err != nil {
		log.Printf("Server init: %v\n", err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
	auth.RegisterAuthorizationServiceServer(grpcServer, &AuthorizationServer{
		service: service,
	})

	if err = grpcServer.Serve(lis); err != nil {
		log.Printf("Server init: %v\n", err)
		os.Exit(1)
	}
}
