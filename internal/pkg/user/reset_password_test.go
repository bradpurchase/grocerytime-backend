package user

import (
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *Suite) TestResetPasswordEmail_TokenExpired() {
	token := "abc123"
	s.mock.ExpectQuery("^SELECT (.+) FROM \"users\"*").
		WithArgs(token).
		WillReturnRows(sqlmock.NewRows([]string{"password_reset_token_expiry"}))

	_, err := ResetPassword("test123", token)
	require.Error(s.T(), err)
	assert.Equal(s.T(), err.Error(), "token invalid or expired")
}

func (s *Suite) TestResetPasswordEmail_TokenValid() {
	token := uuid.NewV4()
	tokenExpiry := time.Now().Add(time.Minute * 5)
	user := models.User{
		ID:                       uuid.NewV4(),
		Email:                    "test@example.com",
		PasswordResetToken:       &token,
		PasswordResetTokenExpiry: &tokenExpiry,
	}
	s.mock.ExpectQuery("^SELECT (.+) FROM \"users\"*").
		WithArgs(token).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow(user.ID, user.Email))

	s.mock.ExpectExec("^UPDATE \"users\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))

	resetUser, err := ResetPassword("test123", token.String())
	require.NoError(s.T(), err)
	assert.Equal(s.T(), resetUser.Email, user.Email)
}
