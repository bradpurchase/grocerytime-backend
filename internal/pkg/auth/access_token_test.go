package auth

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *Suite) TestRetrieveAccessToken_NoAuthHeader() {
	_, err := RetrieveAccessToken("")
	assert.Equal(s.T(), err.Error(), "no authorization header provided")
}

func (s *Suite) TestRetrieveAccessToken_NotBearerToken() {
	_, err := RetrieveAccessToken("Hello123")
	assert.Equal(s.T(), err.Error(), "no bearer token provided")
}

func (s *Suite) TestRetrieveAccessToken_ValidToken() {
	token, err := RetrieveAccessToken("Bearer hello123")
	require.NoError(s.T(), err)
	assert.Equal(s.T(), token, "hello123")
}
