package auth

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestFetchAuthenticatedUser_NoAuthHeader(t *testing.T) {
	dbMock, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db, _ := gorm.Open("postgres", dbMock)

	user, err := FetchAuthenticatedUser(db, "")
	assert.Equal(t, err.Error(), "no authorization header provided")
	assert.Nil(t, user)
}

func TestFetchAuthenticatedUser_MalformedAuthHeader(t *testing.T) {
	dbMock, _, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	user, err := FetchAuthenticatedUser(db, "hello123")
	assert.Equal(t, err.Error(), "no bearer token provided")
	assert.Nil(t, user)
}

func TestFetchAuthenticatedUser_TokenNotFound(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	mock.ExpectQuery("^SELECT (.+) FROM auth_tokens*").
		WithArgs("hello123").
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, e := FetchAuthenticatedUser(db, "Bearer hello123")
	require.Error(t, e)
	assert.Equal(t, e.Error(), "token invalid/expired")
}

func TestFetchAuthenticatedUser_TokenFound(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	testID := "12345678-123a-456b-789c-1a234bcde567"

	authTokenRows := sqlmock.
		NewRows([]string{
			"id",
			"client_id",
			"user_id",
			"access_token",
			"refresh_token",
			"expires_in",
			"created_at",
			"updated_at",
		}).
		AddRow(testID, testID, testID, "hello123", "world456", time.Now().Add(time.Hour*1), time.Now(), time.Now())
	userRows := sqlmock.
		NewRows([]string{
			"id",
			"email",
			"password",
			"first_name",
			"last_name",
			"last_seen_at",
			"created_at",
			"updated_at",
		}).
		AddRow(testID, "test@example.com", "password", "John", "Doe", time.Now(), time.Now(), time.Now())
	mock.ExpectQuery("^SELECT (.+) FROM \"auth_tokens\"*").
		WithArgs("hello123").
		WillReturnRows(authTokenRows)
	mock.ExpectQuery("^SELECT (.+) FROM \"users\"*").
		WithArgs(testID).
		WillReturnRows(userRows)

	user, err := FetchAuthenticatedUser(db, "Bearer hello123")
	require.NoError(t, err)
	assert.NotNil(t, user)
}
