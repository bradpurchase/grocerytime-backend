package user

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *Suite) TestSendForgotPasswordEmail_EmailNotFound() {
	email := "test@example.com"
	s.mock.ExpectQuery("^SELECT (.+) FROM \"users\"*").
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	_, err := SendForgotPasswordEmail(email)
	require.Error(s.T(), err)
	assert.Equal(s.T(), err.Error(), "record not found")
}

func (s *Suite) TestSendForgotPasswordEmail_EmailFound() {
	email := "brad@example.com"
	user := models.User{ID: uuid.NewV4(), Email: email}
	s.mock.ExpectQuery("^SELECT (.+) FROM \"users\"*").
		WithArgs(user.Email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow(user.ID, user.Email))

	s.mock.ExpectExec("^UPDATE \"users\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))

	resetUser, e := SendForgotPasswordEmail(email)
	require.NoError(s.T(), e)
	require.Equal(s.T(), resetUser.Email, email)
}
