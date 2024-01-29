package userServices

import (
	"encoding/base64"
)

type Base64EncodingService struct{}

func NewBase64EncodingService() *Base64EncodingService {
	return &Base64EncodingService{}
}

func (s *Base64EncodingService) GenerateToken(username string) (string, error) {
	return base64.StdEncoding.EncodeToString([]byte(username)), nil
}

func (s *Base64EncodingService) ParseToken(token string) (string, error) {
	response, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return "", err
	}

	return string(response), nil
}
