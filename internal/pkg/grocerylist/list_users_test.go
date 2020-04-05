package grocerylist

import (
	"database/sql/driver"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
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

func TestAddUserToList_UserExistsNotYetAdded(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	listID := uuid.NewV4()
	list := &models.List{ID: listID, Name: "Test List", CreatedAt: time.Now(), UpdatedAt: time.Now()}

	userID := uuid.NewV4()
	email := "test@example.com"
	mock.ExpectQuery("^SELECT (.+) FROM \"users\"*").
		WithArgs(email).
		WillReturnRows(sqlmock.
			NewRows([]string{
				"id",
				"email",
			}).
			AddRow(userID, email))

	mock.ExpectQuery("^SELECT (.+) FROM \"list_users\"*").
		WithArgs(listID, userID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO \"list_users\" (.+)$").
		WithArgs(listID, userID, "", AnyTime{}, AnyTime{}).
		WillReturnRows(sqlmock.NewRows([]string{"list_id"}).AddRow(listID))
	mock.ExpectCommit()

	listUser, err := AddUserToList(db, email, list)
	require.NoError(t, err)
	assert.Equal(t, listUser.(models.ListUser).UserID, userID)
}

func TestAddUserToList_UserExistsAlreadyAdded(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	listID := uuid.NewV4()
	list := &models.List{ID: listID, Name: "Test List", CreatedAt: time.Now(), UpdatedAt: time.Now()}

	userID := uuid.NewV4()
	email := "test@example.com"
	mock.ExpectQuery("^SELECT (.+) FROM \"users\"*").
		WithArgs(email).
		WillReturnRows(sqlmock.
			NewRows([]string{
				"id",
				"email",
			}).
			AddRow(userID, email))

	listUserID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"list_users\"*").
		WithArgs(listID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(listUserID))

	listUser, err := AddUserToList(db, email, list)
	require.NoError(t, err)
	assert.Equal(t, listUser.(models.ListUser).ID, listUserID)
	assert.Equal(t, listUser.(models.ListUser).UserID, userID)
}
