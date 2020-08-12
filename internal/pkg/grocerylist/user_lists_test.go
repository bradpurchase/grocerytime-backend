package grocerylist

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRetrieveUserLists_NoLists(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	user := models.User{ID: uuid.NewV4(), Email: "test@example.com"}
	mock.ExpectQuery("^SELECT (.+) FROM \"lists\"*").
		WithArgs(user.ID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	userLists, err := RetrieveUserLists(db, user)
	require.NoError(t, err)
	assert.Equal(t, userLists, []models.List{})
}

func TestRetrieveUserLists_HasListsCreated(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	user := models.User{ID: uuid.NewV4(), Email: "test@example.com"}
	listRows := sqlmock.
		NewRows([]string{
			"id",
			"user_id",
			"name",
			"created_at",
			"updated_at",
			"deleted_at",
		}).
		AddRow(uuid.NewV4(), user.ID, "Provigo Grocery List", time.Now(), time.Now(), nil).
		AddRow(uuid.NewV4(), user.ID, "Beer Store List", time.Now(), time.Now(), nil)
	mock.
		ExpectQuery("^SELECT lists.* FROM \"lists\"*").
		WithArgs(user.ID).
		WillReturnRows(listRows)

	userLists, err := RetrieveUserLists(db, user)
	require.NoError(t, err)
	assert.Len(t, userLists, 2)
	assert.Equal(t, userLists[0].Name, "Provigo Grocery List")
	assert.Equal(t, userLists[0].UserID, user.ID)
	assert.Equal(t, userLists[1].Name, "Beer Store List")
	assert.Equal(t, userLists[1].UserID, user.ID)
}

func TestRetrieveUserLists_HasListsCreatedAndJoined(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	user := models.User{ID: uuid.NewV4(), Email: "test@example.com"}
	sharingUserID := uuid.NewV4()
	sharedListID := uuid.NewV4()
	listRows := sqlmock.
		NewRows([]string{
			"id",
			"user_id",
			"name",
			"created_at",
			"updated_at",
		}).
		AddRow(uuid.NewV4(), user.ID, "Provigo Grocery List", time.Now(), time.Now()).
		AddRow(uuid.NewV4(), user.ID, "Beer Store List", time.Now(), time.Now()).
		AddRow(sharedListID, sharingUserID, "Shared List", time.Now(), time.Now())
	mock.
		ExpectQuery("^SELECT lists.* FROM \"lists\"*").
		WithArgs(user.ID).
		WillReturnRows(listRows)

	userLists, err := RetrieveUserLists(db, user)
	require.NoError(t, err)
	assert.Len(t, userLists, 3)
	assert.Equal(t, userLists[0].Name, "Provigo Grocery List")
	assert.Equal(t, userLists[0].UserID, user.ID)
	assert.Equal(t, userLists[1].Name, "Beer Store List")
	assert.Equal(t, userLists[2].Name, "Shared List")
	assert.Equal(t, userLists[2].UserID, sharingUserID)
}

func TestRetrieveListForUser_NotFound(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	listID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"lists\"*").
		WithArgs(listID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, e := RetrieveListForUser(db, listID, uuid.NewV4())
	require.Error(t, e)
	assert.Equal(t, e.Error(), "record not found")
}

func TestRetrieveListForUser_ListCreatedByUser(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	listID := uuid.NewV4()
	userID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"lists\"*").
		WithArgs(listID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(listID, "Example List"))

	mock.ExpectQuery("^SELECT (.+) FROM \"list_users\"*").
		WithArgs(listID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"list_id", "user_id"}).AddRow(listID, userID))

	list, err := RetrieveListForUser(db, listID, userID)
	require.NoError(t, err)
	assert.Equal(t, list.Name, "Example List")
}

func TestRetrieveListForUser_ListSharedToUser(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	listID := uuid.NewV4()
	userID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"lists\"*").
		WithArgs(listID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(listID, "Example List"))

	listUserRows := sqlmock.
		NewRows([]string{"list_id", "user_id", "active", "creator"}).
		AddRow(listID, uuid.NewV4(), true, true).
		AddRow(listID, userID, true, false)
	mock.ExpectQuery("^SELECT (.+) FROM \"list_users\"*").
		WithArgs(listID, userID).
		WillReturnRows(listUserRows)

	list, err := RetrieveListForUser(db, listID, userID)
	require.NoError(t, err)
	assert.Equal(t, list.Name, "Example List")
}

func TestRetrieveListForUserByName_NotFound(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	listName := "Example List"
	userID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"lists\"*").
		WithArgs(listName, userID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, e := RetrieveListForUserByName(db, listName, userID)
	require.Error(t, e)
	assert.Equal(t, e.Error(), "record not found")
}

func TestRetrieveListForUserByName_Found(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	listName := "Example List"
	userID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"lists\"*").
		WithArgs(listName, userID).
		WillReturnRows(sqlmock.NewRows([]string{"name", "user_id"}).AddRow(listName, userID))

	list, err := RetrieveListForUserByName(db, listName, userID)
	require.NoError(t, err)
	assert.Equal(t, list.Name, listName)
}

func TestRetrieveSharableList_NotFound(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	listID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"lists\"*").
		WithArgs(listID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, e := RetrieveSharableList(db, listID)
	require.Error(t, e)
	//assert.Equal(t, list.(models.List).Name, listID)
}

func TestRetrieveSharableList_Found(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	listID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"lists\"*").
		WithArgs(listID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(listID, "My Test List"))

	list, e := RetrieveSharableList(db, listID)
	require.NoError(t, e)
	assert.Equal(t, list.ID, listID)
	assert.Equal(t, list.Name, "My Test List")
}

func TestDeleteList_ListNotFound(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	listID := uuid.NewV4()
	userID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"lists\"*").
		WithArgs(listID, userID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, e := DeleteList(db, listID, userID)
	require.Error(t, e)
	assert.Equal(t, e.Error(), "couldn't retrieve list")
}

func TestDeleteList_ListFound(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	listID := uuid.NewV4()
	userID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"lists\"*").
		WithArgs(listID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(listID, userID))

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE \"lists\" (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	tripID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(listID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tripID))

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE \"items\"*").
		WithArgs(AnyTime{}, tripID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE \"grocery_trips\" (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectQuery("^SELECT (.+) FROM \"list_users\"*").
		WithArgs(listID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(listID))

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE \"list_users\"*").
		WithArgs(AnyTime{}, listID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	list, err := DeleteList(db, listID, userID)
	require.NoError(t, err)
	assert.Equal(t, list.ID, listID)
}
