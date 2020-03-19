package auth

import (
	"errors"
	"strings"
)

// RetrieveClientCredentials finds and returns an ApiClient record with the key/secret provided
func RetrieveClientCredentials(authHeader string) ([]string, error) {
	if len(authHeader) == 0 {
		return nil, errors.New("no authorization header provided")
	}
	creds := strings.Split(authHeader, ":")
	if len(creds) != 2 {
		return nil, errors.New("authorization header value is malformed, must be key:secret")
	}
	return creds, nil
}
