package auth

import (
	"context"
	"errors"
	userDomain "github.com/strikersk/user-auth/src/domain"
	userGrpc "github.com/strikersk/user-auth/src/proto/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"log"
)

type AuthorizationClientService struct {
	client userGrpc.AuthorizationServiceClient
}

func NewAuthorizationClientService() *AuthorizationClientService {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v\n", err)
	}

	return &AuthorizationClientService{
		client: userGrpc.NewAuthorizationServiceClient(conn),
	}
}

func (c *AuthorizationClientService) RegisterUser(ctx context.Context, user userDomain.User) error {
	req := &userGrpc.RegisterRequest{
		Username: user.Username,
		Password: user.Password,
	}

	_, err := c.client.RegisterUser(ctx, req)
	return err
}

func (c *AuthorizationClientService) LoginUser(ctx context.Context, user userDomain.User) (string, error) {
	req := &userGrpc.LoginRequest{
		Username: user.Username,
		Password: user.Password,
	}

	response, err := c.client.LoginUser(ctx, req)
	if err != nil {
		msg, _ := status.FromError(err)
		return "", errors.New(msg.Message())
	}

	return response.Token, nil
}
