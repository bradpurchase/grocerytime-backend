package auth

import (
	"database/sql"
	"database/sql/driver"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type Suite struct {
	suite.Suite

	DB   *gorm.DB
	mock sqlmock.Sqlmock
}

func (s *Suite) SetupSuite() {
	var (
		dbMock *sql.DB
		err    error
	)

	dbMock, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)
	s.DB, err = gorm.Open(postgres.New(postgres.Config{Conn: dbMock}), &gorm.Config{})
	require.NoError(s.T(), err)

	db.Manager = s.DB
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestFetchAuthenticatedUser_NoAuthHeader() {
	_, err := FetchAuthenticatedUser("")
	require.Error(s.T(), err)
	assert.Equal(s.T(), err.Error(), "no authorization header provided")
}

func (s *Suite) TestFetchAuthenticatedUser_MalformedAuthHeader() {
	_, err := FetchAuthenticatedUser("hello123")
	require.Error(s.T(), err)
	assert.Equal(s.T(), err.Error(), "no bearer token provided")
}

func (s *Suite) TestFetchAuthenticatedUser_TokenNotFound() {
	s.mock.ExpectQuery("^SELECT (.+) FROM auth_tokens*").
		WithArgs("hello123").
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, e := FetchAuthenticatedUser("Bearer hello123")
	require.Error(s.T(), e)
	assert.Equal(s.T(), e.Error(), "token invalid/expired")
}

func (s *Suite) TestFetchAuthenticatedUser_TokenFound() {
	testID := uuid.NewV4()
	authTokenRows := s.mock.
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
	userRows := s.mock.
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
	s.mock.ExpectQuery("^SELECT (.+) FROM \"auth_tokens\"*").
		WithArgs("hello123").
		WillReturnRows(authTokenRows)
	s.mock.ExpectQuery("^SELECT (.+) FROM \"users\"*").
		WithArgs(testID).
		WillReturnRows(userRows)

	user, err := FetchAuthenticatedUser("Bearer hello123")
	require.NoError(s.T(), err)
	assert.NotNil(s.T(), user)
}
