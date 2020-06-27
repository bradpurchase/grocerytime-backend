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

func TestInviteToListByEmail_UserExistsNotYetAdded(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	listID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"lists\"*").
		WithArgs(listID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(listID))

	email := "test@example.com"
	mock.ExpectQuery("^SELECT (.+) FROM \"list_users\"*").
		WithArgs(listID, email).
		WillReturnRows(sqlmock.NewRows([]string{}))

	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO \"list_users\" (.+)$").
		WithArgs(listID, "00000000-0000-0000-0000-000000000000", email, AnyTime{}, AnyTime{}).
		WillReturnRows(sqlmock.NewRows([]string{"list_id"}).AddRow(listID))
	mock.ExpectQuery("^SELECT name FROM \"lists\"*").
		WithArgs(listID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(listID))
	mock.ExpectCommit()

	listUser, err := InviteToListByEmail(db, listID, email)
	require.NoError(t, err)
	assert.Equal(t, listUser.Email, email)
}

func TestInviteToListByEmail_UserExistsAlreadyAdded(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	listID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"lists\"*").
		WithArgs(listID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(listID))

	email := "test@example.com"
	listUserID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"list_users\"*").
		WithArgs(listID, email).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(listUserID))

	listUser, err := InviteToListByEmail(db, listID, email)
	require.NoError(t, err)
	assert.Equal(t, listUser.ID, listUserID)
	assert.Equal(t, listUser.Email, email)
}

func TestRemoveUserFromList_ListNotFound(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	listID := uuid.NewV4()
	user := models.User{ID: uuid.NewV4()}
	mock.ExpectQuery("^SELECT (.+) FROM \"lists\"*").
		WithArgs(listID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, e := RemoveUserFromList(db, user, listID)
	require.Error(t, e)
	assert.Equal(t, e.Error(), "list not found")
}

func TestRemoveUserFromList_ListUserNotFound(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	listID := uuid.NewV4()
	user := models.User{ID: uuid.NewV4()}
	mock.ExpectQuery("^SELECT (.+) FROM \"lists\"*").
		WithArgs(listID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(listID))

	_, e := RemoveUserFromList(db, user, listID)
	require.Error(t, e)
	assert.Equal(t, e.Error(), "list user not found")
}

func TestRemoveUserFromList_SuccessInvitedUser(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	listID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"lists\"*").
		WithArgs(listID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(listID))

	user := models.User{ID: uuid.NewV4(), Email: "test@example.com"}
	listUser := models.ListUser{ID: uuid.NewV4(), Email: user.Email}
	mock.ExpectQuery("^SELECT (.+) FROM \"list_users\"*").
		WithArgs(listID, user.ID, user.Email).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(listUser.ID))

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM \"list_users\"*").
		WithArgs(listUser.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Test querying for data to send the email about invite being declined
	creatorUser := models.User{ID: uuid.NewV4()}
	creatorListUser := models.ListUser{ID: uuid.NewV4(), UserID: creatorUser.ID}
	mock.ExpectQuery("^SELECT (.+) FROM \"list_users\"*").
		WithArgs(true).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(creatorListUser.ID))
	mock.ExpectQuery("^SELECT (.+) FROM \"users\"*").
		WithArgs(creatorListUser.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(user.ID))

	lu, err := RemoveUserFromList(db, user, listID)
	require.NoError(t, err)
	assert.Equal(t, lu.(*models.ListUser).ID, listUser.ID)
}

func TestAddUserToList_UserNotFoundInList(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	listID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"lists\"*").
		WithArgs(listID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(listID))

	user := models.User{ID: uuid.NewV4(), Email: "test@example.com"}
	listUserID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"list_users\"*").
		WithArgs(listID, user.Email).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(listUserID))

	_, e := AddUserToList(db, user, listID)
	require.Error(t, e)
}

func TestAddUserToList_Success(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	listID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"lists\"*").
		WithArgs(listID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(listID))

	email := "test@example.com"
	user := models.User{ID: uuid.NewV4(), Email: email}
	listUserID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"list_users\"*").
		WithArgs(listID, user.Email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow(listUserID, email))

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE \"list_users\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	listUser, err := AddUserToList(db, user, listID)
	require.NoError(t, err)
	assert.Equal(t, listUser.(*models.ListUser).ID, listUserID)
	assert.Equal(t, listUser.(*models.ListUser).UserID, user.ID)
	assert.Equal(t, listUser.(*models.ListUser).Email, "")
}

func TestRetrieveListUsers_HasListUsers(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	listID := uuid.NewV4()
	list := &models.List{ID: listID, Name: "Test List", CreatedAt: time.Now(), UpdatedAt: time.Now()}

	listUserRows := sqlmock.
		NewRows([]string{
			"id",
			"list_id",
			"user_id",
			"creator",
			"active",
			"created_at",
			"updated_at",
		}).
		AddRow(uuid.NewV4(), listID, uuid.NewV4(), true, true, time.Now(), time.Now()).
		AddRow(uuid.NewV4(), listID, uuid.NewV4(), false, true, time.Now(), time.Now())
	mock.ExpectQuery("^SELECT (.+) FROM \"list_users\"*").
		WithArgs(listID).
		WillReturnRows(listUserRows)

	listUsers, err := RetrieveListUsers(db, list.ID)
	require.NoError(t, err)
	assert.Equal(t, len(listUsers.([]models.ListUser)), 2)
	assert.Equal(t, listUsers.([]models.ListUser)[0].ListID, list.ID)
}
