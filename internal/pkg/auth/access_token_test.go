package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRetrieveAccessToken_NoAuthHeader(t *testing.T) {
	_, err := RetrieveAccessToken("")
	assert.Equal(t, err.Error(), "no authorization header provided")
}

func TestRetrieveAccessToken_NotBearerToken(t *testing.T) {
	_, err := RetrieveAccessToken("Hello123")
	assert.Equal(t, err.Error(), "no bearer token provided")
}

func TestRetrieveAccessToken_ValidToken(t *testing.T) {
	token, err := RetrieveAccessToken("Bearer hello123")
	require.NoError(t, err)
	assert.Equal(t, token, "hello123")
}
