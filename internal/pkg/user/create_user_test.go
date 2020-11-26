package user

import (
	"github.com/DATA-DOG/go-sqlmock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *Suite) TestCreateUser_InvalidPassword() {
	_, e := CreateUser("test@example.com", "", "John", uuid.NewV4())
	require.Error(s.T(), e)
}

func (s *Suite) TestCreateUser_DuplicateEmail() {
	email := "test@example.com"
	s.mock.ExpectQuery("^SELECT (.+) FROM \"users\"*").
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"count(1)"}).AddRow(1))

	_, e := CreateUser(email, "password", "John", uuid.NewV4())
	require.Error(s.T(), e)
	assert.Equal(s.T(), e.Error(), "An account with this email address already exists")
}

func (s *Suite) TestCreateUser_UserCreated() {
	email := "test@example.com"
	s.mock.ExpectQuery("^SELECT (.+) FROM \"users\"*").
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"count(1)"}).AddRow(0))

	userID := uuid.NewV4()
	name := "John Doe"
	s.mock.ExpectQuery("^INSERT INTO \"users\" (.+)$").
		WithArgs(email, sqlmock.AnyArg(), name, nil, nil, AnyTime{}, AnyTime{}, AnyTime{}).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userID))

	s.mock.ExpectExec("^DELETE FROM \"auth_tokens\"*").
		WillReturnResult(sqlmock.NewResult(1, 1))

	clientID := uuid.NewV4()
	s.mock.ExpectQuery("^INSERT INTO \"auth_tokens\" (.+)$").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), AnyTime{}, AnyTime{}, AnyTime{}, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "client_id", "user_id"}).AddRow(uuid.NewV4(), clientID, userID))

	user, err := CreateUser(email, "password", name, clientID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), user.Email, email)
	assert.Equal(s.T(), user.Name, name)
	assert.NotNil(s.T(), user.Tokens)
}
