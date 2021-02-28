package auth

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/lestrrat-go/jwx/jwk"
)

// SignInWithApple will verify an identityToken
func SignInWithApple(identityToken string, email string, name string) (interface{}, error) {
	token, err := jwt.Parse(identityToken, VerifyTokenSignature)
	if err != nil {
		return nil, err
	}
	claims := token.Claims.(jwt.MapClaims)
	for key, val := range claims {
		fmt.Printf("%s\t%v\n", key, val)
	}
	return nil, nil
}

// VerifyTokenSignature fetches Apple's public key for verifying the ID token signature
func VerifyTokenSignature(token *jwt.Token) (interface{}, error) {
	jwksURL := "https://appleid.apple.com/auth/keys"
	keySet, err := jwk.Fetch(context.Background(), jwksURL)
	if err != nil {
		log.Printf("failed to parse JWK: %s", err)
		return nil, err
	}

	kid, ok := token.Header["kid"].(string)
	if !ok {
		return nil, errors.New("expecting JWT header to have string kid")
	}

	key, ok := keySet.LookupKeyID(kid)
	if !ok {
		return nil, fmt.Errorf("unable to find key %q", kid)
	}

	var rawKey interface{}
	if err := key.Raw(&rawKey); err != nil {
		return nil, err
	}
	return rawKey, nil
}
