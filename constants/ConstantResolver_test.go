package constants

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrorTypeResolving(t *testing.T) {
	testError := errors.New("test error")
	actualMessage := resolveString(testError)
	assert.Equal(t, "test error", actualMessage, "Message were not resolved properly")
}

func TestStringTypeResolving(t *testing.T) {
	testError := "test error"
	actualMessage := resolveString(testError)
	assert.Equal(t, "test error", actualMessage, "Message were not resolved properly")
}

func TestUnknownTypeResolving(t *testing.T) {
	testError := 123
	actualMessage := resolveString(testError)
	assert.Equal(t, "value could not be resolved", actualMessage, "Message were not resolved properly")
}
