package userServices

import (
	"encoding/base64"
	"github.com/strikersk/user-auth/src/domain"
)

type Base64EncodingService struct{}

func NewBase64EncodingService() Base64EncodingService {
	return Base64EncodingService{}
}

func (s Base64EncodingService) GenerateToken(user domain.UserDTO) (string, error) {
	return base64.URLEncoding.EncodeToString([]byte(user.Username)), nil
}

func (s Base64EncodingService) ParseToken(token string) (string, error) {
	response, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return "", err
	}

	return string(response), nil
}
