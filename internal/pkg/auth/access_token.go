package auth

import (
	"errors"
	"strings"
)

// RetrieveAccessToken finds the access token within the Authorization header string
func RetrieveAccessToken(authHeader string) (string, error) {
	if len(authHeader) == 0 {
		return authHeader, errors.New("no authorization header provided")
	}
	bearer := strings.Split(authHeader, "Bearer ")
	if len(bearer) != 2 {
		return authHeader, errors.New("no bearer token provided")
	}
	token := bearer[1]
	return token, nil
}
