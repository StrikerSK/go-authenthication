package auth

import (
	"context"
	"github.com/Pallinder/go-randomdata"
	"github.com/strikersk/user-auth/src/domain"
	"log"
)

type AuthorizationSample struct {
	client *AuthorizationClientService
}

func NewAuthorizationSample() *AuthorizationSample {
	return &AuthorizationSample{
		client: NewAuthorizationClientService(),
	}
}

func (r *AuthorizationSample) RegisterUser() {
	randomProfile := randomdata.GenerateProfile(randomdata.RandomGender)
	req := domain.User{
		UserCredentials: domain.UserCredentials{
			Username: randomProfile.Login.Username,
			Password: randomProfile.Login.Password,
		},
	}

	log.Println("Attempting to register user")
	if err := r.client.RegisterUser(context.Background(), req); err != nil {
		log.Println(err)
	} else {
		log.Println("User registered successfully")
	}
}

func (r *AuthorizationSample) LoginUser() {
	req := domain.User{
		UserCredentials: domain.UserCredentials{
			Username: "test",
			Password: "test",
		},
	}

	log.Println("Attempting to login user")
	token, err := r.client.LoginUser(context.Background(), req)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Token returned:", token)
	}
}
