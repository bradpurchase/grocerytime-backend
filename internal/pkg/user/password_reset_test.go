package user

import (
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *Suite) TestVerifyPasswordResetToken_TokenExpired() {
	token := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"users\"*").
		WithArgs(token).
		WillReturnRows(sqlmock.NewRows([]string{"password_reset_token_expiry"}))

	_, err := VerifyPasswordResetToken(token)
	require.Error(s.T(), err)
	assert.Equal(s.T(), err.Error(), "token expired")
}

func (s *Suite) TestVerifyPasswordResetToken_TokenValid() {
	token := uuid.NewV4()
	tokenExpiry := time.Now().Add(time.Minute * 5)
	s.mock.ExpectQuery("^SELECT (.+) FROM \"users\"*").
		WithArgs(token).
		WillReturnRows(sqlmock.NewRows([]string{"password_reset_token_expiry"}).AddRow(tokenExpiry))

	_, err := VerifyPasswordResetToken(token)
	require.NoError(s.T(), err)
}
