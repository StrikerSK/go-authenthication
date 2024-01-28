package userServices

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/strikersk/user-auth/src/domain"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestPasswordEncryption(t *testing.T) {
	user := domain.UserCredentials{
		Username: "admin",
		Password: "admin",
	}
	passwordService := BcryptUserPasswordService{
		cost: bcrypt.MinCost,
	}
	err := passwordService.SetPassword(&user)

	fmt.Println(user.Password)
	assert.Nil(t, err, "Error should not be returned")
	assert.NotEqual(t, "admin", user.Password, "Received password: "+user.Password)
}

func TestPasswordValidation(t *testing.T) {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.MinCost)

	assertedUser := domain.UserCredentials{
		Username: "admin",
		Password: string(encryptedPassword),
	}

	expectedUser := domain.UserCredentials{
		Username: "admin",
		Password: "admin",
	}

	passwordService := BcryptUserPasswordService{}
	err = passwordService.ValidatePassword(assertedUser, expectedUser)

	assert.Nil(t, err, "Error should not be returned")
}
