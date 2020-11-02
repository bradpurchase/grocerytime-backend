package auth

import (
	"github.com/stretchr/testify/assert"
)

func (s *Suite) TestRetrieveClientCredentials_NoAuthHeader() {
	_, err := RetrieveClientCredentials("")
	assert.Equal(s.T(), err.Error(), "no authorization header provided")
}

func (s *Suite) TestRetrieveClientCredentials_Malformed() {
	_, err := RetrieveClientCredentials("Hello123")
	assert.Equal(s.T(), err.Error(), "authorization header value is malformed, must be key:secret")
}

func (s *Suite) TestRetrieveClientCredentials() {
	token, _ := RetrieveClientCredentials("hello123:world456")
	assert.Equal(s.T(), token, []string{"hello123", "world456"})
}
