package stores

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestRetrieveUserStores_NoStores(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: dbMock}), &gorm.Config{})
	require.NoError(t, err)

	user := models.User{ID: uuid.NewV4(), Email: "test@example.com"}
	mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(user.ID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	userStores, err := RetrieveUserStores(db, user)
	require.NoError(t, err)
	assert.Equal(t, userStores, []models.Store{})
}

func TestRetrieveUserStores_HasStoresCreated(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: dbMock}), &gorm.Config{})
	require.NoError(t, err)

	user := models.User{ID: uuid.NewV4(), Email: "test@example.com"}
	storeRows := sqlmock.
		NewRows([]string{
			"id",
			"user_id",
			"name",
			"created_at",
			"updated_at",
			"deleted_at",
		}).
		AddRow(uuid.NewV4(), user.ID, "Provigo Grocery Store", time.Now(), time.Now(), nil).
		AddRow(uuid.NewV4(), user.ID, "Beer Store Store", time.Now(), time.Now(), nil)
	mock.
		ExpectQuery("^SELECT stores.* FROM \"stores\"*").
		WithArgs(user.ID).
		WillReturnRows(storeRows)

	userStores, err := RetrieveUserStores(db, user)
	require.NoError(t, err)
	assert.Len(t, userStores, 2)
	assert.Equal(t, userStores[0].Name, "Provigo Grocery Store")
	assert.Equal(t, userStores[0].UserID, user.ID)
	assert.Equal(t, userStores[1].Name, "Beer Store Store")
	assert.Equal(t, userStores[1].UserID, user.ID)
}

func TestRetrieveUserStores_HasStoresCreatedAndJoined(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: dbMock}), &gorm.Config{})
	require.NoError(t, err)

	user := models.User{ID: uuid.NewV4(), Email: "test@example.com"}
	sharingUserID := uuid.NewV4()
	sharedStoreID := uuid.NewV4()
	storeRows := sqlmock.
		NewRows([]string{
			"id",
			"user_id",
			"name",
			"created_at",
			"updated_at",
		}).
		AddRow(uuid.NewV4(), user.ID, "Provigo Grocery Store", time.Now(), time.Now()).
		AddRow(uuid.NewV4(), user.ID, "Beer Store Store", time.Now(), time.Now()).
		AddRow(sharedStoreID, sharingUserID, "Shared Store", time.Now(), time.Now())
	mock.
		ExpectQuery("^SELECT stores.* FROM \"stores\"*").
		WithArgs(user.ID).
		WillReturnRows(storeRows)

	userStores, err := RetrieveUserStores(db, user)
	require.NoError(t, err)
	assert.Len(t, userStores, 3)
	assert.Equal(t, userStores[0].Name, "Provigo Grocery Store")
	assert.Equal(t, userStores[0].UserID, user.ID)
	assert.Equal(t, userStores[1].Name, "Beer Store Store")
	assert.Equal(t, userStores[2].Name, "Shared Store")
	assert.Equal(t, userStores[2].UserID, sharingUserID)
}

func TestRetrieveInvitedUserStores_NoneFound(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: dbMock}), &gorm.Config{})
	require.NoError(t, err)

	user := models.User{ID: uuid.NewV4(), Email: "test@example.com"}
	mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(user.Email, false).
		WillReturnRows(sqlmock.NewRows([]string{}))

	userStores, err := RetrieveInvitedUserStores(db, user)
	require.NoError(t, err)
	assert.Equal(t, userStores, []models.Store{})
}

func TestRetrieveInvitedUserStores_ResultsFound(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: dbMock}), &gorm.Config{})
	require.NoError(t, err)

	user := models.User{ID: uuid.NewV4(), Email: "test@example.com"}
	storeRows := sqlmock.
		NewRows([]string{
			"id",
			"user_id",
			"name",
			"created_at",
			"updated_at",
			"deleted_at",
		}).
		AddRow(uuid.NewV4(), user.ID, "Provigo Grocery Store", time.Now(), time.Now(), nil).
		AddRow(uuid.NewV4(), user.ID, "Beer Store Store", time.Now(), time.Now(), nil)
	mock.
		ExpectQuery("^SELECT stores.* FROM \"stores\"*").
		WithArgs(user.Email, false).
		WillReturnRows(storeRows)

	userStores, err := RetrieveInvitedUserStores(db, user)
	require.NoError(t, err)
	assert.Len(t, userStores, 2)
	assert.Equal(t, userStores[0].Name, "Provigo Grocery Store")
	assert.Equal(t, userStores[0].UserID, user.ID)
	assert.Equal(t, userStores[1].Name, "Beer Store Store")
	assert.Equal(t, userStores[1].UserID, user.ID)
}

func TestRetrieveStoreForUser_NotFound(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: dbMock}), &gorm.Config{})
	require.NoError(t, err)

	storeID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, e := RetrieveStoreForUser(db, storeID, uuid.NewV4())
	require.Error(t, e)
	assert.Equal(t, e.Error(), "record not found")
}

func TestRetrieveStoreForUser_StoreCreatedByUser(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: dbMock}), &gorm.Config{})
	require.NoError(t, err)

	storeID := uuid.NewV4()
	userID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(storeID, "Example Store"))

	mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"store_id", "user_id"}).AddRow(storeID, userID))

	store, err := RetrieveStoreForUser(db, storeID, userID)
	require.NoError(t, err)
	assert.Equal(t, store.Name, "Example Store")
}

func TestRetrieveStoreForUser_StoreSharedToUser(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: dbMock}), &gorm.Config{})
	require.NoError(t, err)

	storeID := uuid.NewV4()
	userID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(storeID, "Example Store"))

	storeUserRows := sqlmock.
		NewRows([]string{"store_id", "user_id", "active", "creator"}).
		AddRow(storeID, uuid.NewV4(), true, true).
		AddRow(storeID, userID, true, false)
	mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID, userID).
		WillReturnRows(storeUserRows)

	store, err := RetrieveStoreForUser(db, storeID, userID)
	require.NoError(t, err)
	assert.Equal(t, store.Name, "Example Store")
}

func TestRetrieveStoreForUserByName_NotFound(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: dbMock}), &gorm.Config{})
	require.NoError(t, err)

	storeName := "Example Store"
	userID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeName, userID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, e := RetrieveStoreForUserByName(db, storeName, userID)
	require.Error(t, e)
	assert.Equal(t, e.Error(), "record not found")
}

func TestRetrieveStoreForUserByName_Found(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: dbMock}), &gorm.Config{})
	require.NoError(t, err)

	storeName := "Example Store"
	userID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeName, userID).
		WillReturnRows(sqlmock.NewRows([]string{"name", "user_id"}).AddRow(storeName, userID))

	store, err := RetrieveStoreForUserByName(db, storeName, userID)
	require.NoError(t, err)
	assert.Equal(t, store.Name, storeName)
}

func TestDeleteStore_StoreNotFound(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: dbMock}), &gorm.Config{})
	require.NoError(t, err)

	storeID := uuid.NewV4()
	userID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeID, userID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, e := DeleteStore(db, storeID, userID)
	require.Error(t, e)
	assert.Equal(t, e.Error(), "couldn't retrieve store")
}

func TestDeleteStore_StoreFound(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: dbMock}), &gorm.Config{})
	require.NoError(t, err)

	storeID := uuid.NewV4()
	userID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(storeID, userID))

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE \"stores\" (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	tripID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(storeID).
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

	mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(storeID))

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE \"store_users\"*").
		WithArgs(AnyTime{}, storeID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	store, err := DeleteStore(db, storeID, userID)
	require.NoError(t, err)
	assert.Equal(t, store.ID, storeID)
}
