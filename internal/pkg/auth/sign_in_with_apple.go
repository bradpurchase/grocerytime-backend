package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

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

	if err := VerifyIdentityToken(token, nonce, appScheme); err != nil {
		return nil, err
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

// VerifyIdentityToken verifies the identity token following the creteria specified by Apple
// see: https://developer.apple.com/documentation/sign_in_with_apple/sign_in_with_apple_rest_api/verifying_a_user
func VerifyIdentityToken(token *jwt.Token, nonce string, appScheme string) (err error) {
	claims := token.Claims.(jwt.MapClaims)
	for key, val := range claims {
		fmt.Printf("%s\t%v\n", key, val)

		switch field := key; field {
		case "nonce":
			nonceClaim := val.(string)
			if err := VerifyNonce(nonceClaim, nonce); err != nil {
				return err
			}
		case "iss":
			iss := val.(string)
			if err := VerifyIss(iss); err != nil {
				return err
			}
		case "aud":
			aud := val.(string)
			if err := VerifyAud(aud, appScheme); err != nil {
				return err
			}
		case "exp":
			exp := val.(int64)
			if err := VerifyExp(exp); err != nil {
				return err
			}
		}
	}
	return
}

// VerifyNonce verifies that there is a match between the nonce in the JWT claims
// and the nonce value passed down to the server from the SIWA request
func VerifyNonce(nonceClaim, nonceValue string) (err error) {
	if nonceClaim != nonceValue {
		return errors.New("invalid signin (nonce mismatch)")
	}
	return
}

// VerifyIss verifies that the iss field in the claims contains https://appleid.apple.com
func VerifyIss(iss string) (err error) {
	if !strings.Contains(iss, issAppleID) {
		return errors.New("invalid signin (invalid iss)")
	}
	return
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
	return
}

// VerifyExp verifies that the exp field, repesenting expiration time, has not passed
func VerifyExp(exp int64) (err error) {
	expTime := time.Unix(exp, 0)
	time := time.Now()
	if time.After(expTime) {
		return errors.New("invalid signin (identity token expired)")
	}
	return
}
