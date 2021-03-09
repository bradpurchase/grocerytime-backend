package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/mailer"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/lestrrat-go/jwx/jwk"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	jwksURL    = "https://appleid.apple.com/auth/keys"
	issAppleID = "https://appleid.apple.com"
	bundleID   = os.Getenv("APP_BUNDLE_IDENTIFIER")
)

// SignInWithApple will verify an identityToken
func SignInWithApple(identityToken, nonce, email, name, appScheme string, clientID uuid.UUID) (user *models.User, err error) {
	token, err := jwt.Parse(identityToken, VerifyTokenSignature)
	if err != nil {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)
	if err := VerifyIdentityToken(claims, nonce, appScheme); err != nil {
		return nil, err
	}

	user, err = FindOrCreateUserFromIdentityToken(claims, name, clientID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// VerifyTokenSignature fetches Apple's public key for verifying the ID token signature
// see: https://stackoverflow.com/questions/41077953/go-language-and-verify-jwt
func VerifyTokenSignature(token *jwt.Token) (interface{}, error) {
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
func VerifyIdentityToken(claims map[string]interface{}, nonce string, appScheme string) (err error) {
	for key, val := range claims {
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
			exp := val.(float64)
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
	//
	// Note: we check if the scheme is already contained in the bundleID because
	// if you try two SIWA attempts in quick succession, the bundleID already has the scheme
	// (not sure if this would happen in practice, need more testing to verify)
	appScheme = strings.ToLower(appScheme)
	fmt.Printf("[SIWA/VerifyAud] appScheme: %v\n", appScheme)
	if appScheme != "release" && !strings.Contains(bundleID, appScheme) {
		bundleID = fmt.Sprintf("%v.%v", bundleID, appScheme)
	}
	fmt.Printf("[SIWA/VerifyAud] aud: %v/ bundleID: %v\n", aud, bundleID)
	if aud != bundleID {
		return errors.New("invalid signin (aud mismatch)")
	}
	return
}

// VerifyExp verifies that the exp field, repesenting expiration time, has not passed
func VerifyExp(exp float64) (err error) {
	expTime := time.Unix(int64(exp), 0)
	time := time.Now()
	if time.After(expTime) {
		return errors.New("invalid signin (identity token expired)")
	}
	return
}

// FindOrCreateUserFromIdentityToken finds or creates a user from the identity token
func FindOrCreateUserFromIdentityToken(claims map[string]interface{}, userName string, clientID uuid.UUID) (user *models.User, err error) {
	// Check if there's a user that matches the sub (siwa_id) or email included in the token claims
	sub := claims["sub"].(string)
	email := claims["email"].(string)
	if user := FindUserBySubOrEmail(sub, email); user != nil {
		return user, nil
	}

	// If no user was found, we create one and associate the siwa_id for further logins
	if user, err = CreateUserFromIdentityToken(sub, userName, email, clientID); err != nil {
		return nil, err
	}

	// Create an access token on our side
	authToken := &models.AuthToken{
		UserID:     user.ID,
		ClientID:   clientID,
		DeviceName: "SIWA",
	}
	if err := db.Manager.Create(&authToken).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// FindUserBySubOrEmail tries to find a user with the sub (ID from Apple's side) or email
func FindUserBySubOrEmail(sub, email string) (user *models.User) {
	var foundUser models.User
	if err := db.Manager.Where("siwa_id = ?", sub).Or("email = ?", email).First(&foundUser).Error; err != nil {
		return nil
	}

	// Update the user to associate siwa_id
	updateUserQuery := db.Manager.
		Model(&foundUser).
		Where("id = ?", foundUser.ID).
		Updates(&models.User{SiwaID: &sub}).
		Error
	if err := updateUserQuery; err != nil {
		return nil
	}
	return &foundUser
}

// CreateUserFromIdentityToken creates a user from identity token claims
func CreateUserFromIdentityToken(sub, userName, email string, clientID uuid.UUID) (user *models.User, err error) {
	password := utils.RandString(16) // fake password to persist the user
	passhash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user = &models.User{
		Name:       userName,
		Email:      email,
		Password:   string(passhash),
		LastSeenAt: time.Now(),
		SiwaID:     &sub,
	}
	if err := db.Manager.Create(&user).Error; err != nil {
		return nil, err
	}

	// Send an email upon user creation
	_, mailErr := mailer.SendNewUserEmail(email)
	if mailErr != nil {
		return nil, mailErr
	}

	return user, nil
}
