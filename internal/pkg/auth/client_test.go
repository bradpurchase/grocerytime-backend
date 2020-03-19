package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRetrieveClientCredentials_NoAuthHeader(t *testing.T) {
	_, err := RetrieveClientCredentials("")
	assert.Equal(t, err.Error(), "no authorization header provided")
}

func TestRetrieveClientCredentials_Malformed(t *testing.T) {
	_, err := RetrieveClientCredentials("Hello123")
	assert.Equal(t, err.Error(), "authorization header value is malformed, must be key:secret")
}

func TestRetrieveClientCredentials(t *testing.T) {
	token, _ := RetrieveClientCredentials("hello123:world456")
	assert.Equal(t, token, []string{"hello123", "world456"})
}
