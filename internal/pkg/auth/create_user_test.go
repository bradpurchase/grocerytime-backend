package auth

import (
	"database/sql/driver"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestCreateUser_InvalidPassword(t *testing.T) {
	dbMock, _, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	email := "test@example.com"
	_, e := CreateUser(db, email, "", uuid.NewV4())
	require.Error(t, e)
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	email := "test@example.com"
	mock.ExpectQuery("^SELECT (.+) FROM \"users\"*").
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"email"}).AddRow(email))

	_, e := CreateUser(db, email, "password", uuid.NewV4())
	require.Error(t, e)
	assert.Equal(t, e.Error(), "An account with this email address already exists")
}

func TestCreateUser_UserCreated(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	email := "test@example.com"
	listName := "My Grocery List"

	mock.ExpectQuery("^SELECT (.+) FROM \"users\"*").
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{}))

	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO \"users\" (.+)$").
		WithArgs(email, sqlmock.AnyArg(), "", "", AnyTime{}, AnyTime{}, AnyTime{}).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))
	mock.ExpectQuery("^INSERT INTO \"lists\" (.+)$").
		WithArgs(sqlmock.AnyArg(), listName, AnyTime{}, AnyTime{}).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))
	mock.ExpectQuery("^INSERT INTO \"list_users\" (.+)$").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "", true, AnyTime{}, AnyTime{}).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))
	mock.ExpectQuery("^INSERT INTO \"grocery_trips\" (.+)$").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), AnyTime{}, AnyTime{}).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))
	mock.ExpectQuery("^INSERT INTO \"auth_tokens\" (.+)$").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), AnyTime{}, AnyTime{}, AnyTime{}).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))
	mock.ExpectCommit()

	clientID := uuid.NewV4()
	user, err := CreateUser(db, email, "password", clientID)
	require.NoError(t, err)
	assert.Equal(t, user.Email, email)
	assert.Equal(t, user.Lists[0].Name, listName)
	assert.Equal(t, user.Tokens[0].ClientID, clientID)
}
