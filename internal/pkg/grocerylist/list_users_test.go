package grocerylist

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddUserToList_UserDoesntExist(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	email := "test@example.com"
	list := &models.List{Name: "Test List"}

	mock.ExpectQuery("^SELECT (.+) FROM \"users\"*").
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{}))

	userLists, err := AddUserToList(db, email, list)
	require.NoError(t, err)
	assert.Equal(t, userLists, &models.ListUser{})
}

func TestAddUserToList_UserExists(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	email := "test@example.com"
	list := &models.List{Name: "Test List"}
	userID := uuid.NewV4()

	mock.ExpectQuery("^SELECT (.+) FROM \"users\"*").
		WithArgs(email).
		WillReturnRows(sqlmock.
			NewRows([]string{
				"id",
				"email",
			}).
			AddRow(userID, email))

	mock.ExpectQuery("^SELECT (.+) FROM \"list_users\"*").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	//TODOmock.ExpectQuery("INSERT") because this is the case where it creates a new ListUser

	userLists, err := AddUserToList(db, email, list)
	require.NoError(t, err)
	assert.Equal(t, userLists, &models.ListUser{})
}
