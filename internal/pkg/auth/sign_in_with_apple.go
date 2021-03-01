package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/lestrrat-go/jwx/jwk"
)

var (
	issAppleID = "https://appleid.apple.com"
	bundleID   = os.Getenv("APP_BUNDLE_IDENTIFIER")
)

// SignInWithApple will verify an identityToken
func SignInWithApple(identityToken, nonce, email, name, appScheme string) (interface{}, error) {
	token, err := jwt.Parse(identityToken, VerifyTokenSignature)
	if err != nil {
		return nil, err
	}
	claims := token.Claims.(jwt.MapClaims)
	for key, val := range claims {
		fmt.Printf("%s\t%v\n", key, val)

		// Verify the identity token
		// see: https://developer.apple.com/documentation/sign_in_with_apple/sign_in_with_apple_rest_api/verifying_a_user
		switch field := key; field {
		case "nonce":
			nonceClaim := val.(string)
			if err := VerifyNonce(nonceClaim, nonce); err != nil {
				return nil, err
			}
		case "iss":
			iss := val.(string)
			if err := VerifyIss(iss); err != nil {
				return nil, err
			}
		case "aud":
			aud := val.(string)
			if err := VerifyAud(aud, appScheme); err != nil {
				return nil, err
			}
		}
	}
	return nil, nil
}

// VerifyTokenSignature fetches Apple's public key for verifying the ID token signature
// see: https://stackoverflow.com/questions/41077953/go-language-and-verify-jwt
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

// VerifyNonce verifies that there is a match between the nonce in the JWT claims
// and the nonce value passed down to the server from the SIWA request
func VerifyNonce(nonceClaim, nonceValue string) (err error) {
	if nonceClaim != nonceValue {
		return errors.New("invalid signin (nonce mismatch)")
	}
	return nil
}

// VerifyIss verifies that the iss field in the claims contains https://appleid.apple.com
func VerifyIss(iss string) (err error) {
	if !strings.Contains(iss, issAppleID) {
		return errors.New("invalid signin (invalid iss)")
	}
	return nil
}

// VerifyAud verifies that the aud field in the claims matches the app's bundle identifier
func VerifyAud(aud, appScheme string) (err error) {
	// Bundle identifier for non-release schemes will take a different format
	// (e.g. com.supercoolapps.testapp.beta)
	if appScheme != "Release" {
		bundleID = fmt.Sprintf("%v.%v", bundleID, strings.ToLower(appScheme))
	}
	if aud != bundleID {
		return errors.New("invalid signin (aud mismatch)")
	}
	return nil
}
