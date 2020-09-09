package stores

import (
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *Suite) TestRetrieveUserStores_NoStores() {
	user := models.User{ID: uuid.NewV4(), Email: "test@example.com"}
	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(user.ID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	userStores, err := RetrieveUserStores(user)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), userStores, []models.Store{})
}

func (s *Suite) TestRetrieveUserStores_HasStoresCreated() {
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
	s.mock.
		ExpectQuery("^SELECT stores.* FROM \"stores\"*").
		WithArgs(user.ID).
		WillReturnRows(storeRows)

	userStores, err := RetrieveUserStores(user)
	require.NoError(s.T(), err)
	assert.Len(s.T(), userStores, 2)
	assert.Equal(s.T(), userStores[0].Name, "Provigo Grocery Store")
	assert.Equal(s.T(), userStores[0].UserID, user.ID)
	assert.Equal(s.T(), userStores[1].Name, "Beer Store Store")
	assert.Equal(s.T(), userStores[1].UserID, user.ID)
}

func (s *Suite) TestRetrieveUserStores_HasStoresCreatedAndJoined() {
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
	s.mock.
		ExpectQuery("^SELECT stores.* FROM \"stores\"*").
		WithArgs(user.ID).
		WillReturnRows(storeRows)

	userStores, err := RetrieveUserStores(user)
	require.NoError(s.T(), err)
	assert.Len(s.T(), userStores, 3)
	assert.Equal(s.T(), userStores[0].Name, "Provigo Grocery Store")
	assert.Equal(s.T(), userStores[0].UserID, user.ID)
	assert.Equal(s.T(), userStores[1].Name, "Beer Store Store")
	assert.Equal(s.T(), userStores[2].Name, "Shared Store")
	assert.Equal(s.T(), userStores[2].UserID, sharingUserID)
}

func (s *Suite) TestRetrieveInvitedUserStores_NoneFound() {
	user := models.User{ID: uuid.NewV4(), Email: "test@example.com"}
	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(user.Email, false).
		WillReturnRows(sqlmock.NewRows([]string{}))

	userStores, err := RetrieveInvitedUserStores(user)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), userStores, []models.Store{})
}

func (s *Suite) TestRetrieveInvitedUserStores_ResultsFound() {
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
	s.mock.
		ExpectQuery("^SELECT stores.* FROM \"stores\"*").
		WithArgs(user.Email, false).
		WillReturnRows(storeRows)

	userStores, err := RetrieveInvitedUserStores(user)
	require.NoError(s.T(), err)
	assert.Len(s.T(), userStores, 2)
	assert.Equal(s.T(), userStores[0].Name, "Provigo Grocery Store")
	assert.Equal(s.T(), userStores[0].UserID, user.ID)
	assert.Equal(s.T(), userStores[1].Name, "Beer Store Store")
	assert.Equal(s.T(), userStores[1].UserID, user.ID)
}

func (s *Suite) TestRetrieveStoreForUser_NotFound() {
	storeID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, e := RetrieveStoreForUser(storeID, uuid.NewV4())
	require.Error(s.T(), e)
	assert.Equal(s.T(), e.Error(), "record not found")
}

func (s *Suite) TestRetrieveStoreForUser_StoreCreatedByUser() {
	storeID := uuid.NewV4()
	userID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(storeID, "Example Store"))

	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"store_id", "user_id"}).AddRow(storeID, userID))

	store, err := RetrieveStoreForUser(storeID, userID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), store.Name, "Example Store")
}

func (s *Suite) TestRetrieveStoreForUser_StoreSharedToUser() {
	storeID := uuid.NewV4()
	userID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(storeID, "Example Store"))

	storeUserRows := sqlmock.
		NewRows([]string{"store_id", "user_id", "active", "creator"}).
		AddRow(storeID, uuid.NewV4(), true, true).
		AddRow(storeID, userID, true, false)
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID, userID).
		WillReturnRows(storeUserRows)

	store, err := RetrieveStoreForUser(storeID, userID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), store.Name, "Example Store")
}

func (s *Suite) TestRetrieveStoreForUserByName_NotFound() {
	storeName := "Example Store"
	userID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeName, userID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, e := RetrieveStoreForUserByName(storeName, userID)
	require.Error(s.T(), e)
	assert.Equal(s.T(), e.Error(), "record not found")
}

func (s *Suite) TestRetrieveStoreForUserByName_Found() {
	storeName := "Example Store"
	userID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeName, userID).
		WillReturnRows(sqlmock.NewRows([]string{"name", "user_id"}).AddRow(storeName, userID))

	store, err := RetrieveStoreForUserByName(storeName, userID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), store.Name, storeName)
}

func (s *Suite) TestDeleteStore_StoreNotFound() {
	storeID := uuid.NewV4()
	userID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeID, userID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, e := DeleteStore(storeID, userID)
	require.Error(s.T(), e)
	assert.Equal(s.T(), e.Error(), "couldn't retrieve store")
}

func (s *Suite) TestDeleteStore_StoreFound() {
	storeID := uuid.NewV4()
	userID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(storeID, userID))

	s.mock.ExpectBegin()
	s.mock.ExpectExec("^UPDATE \"stores\" (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	tripID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(storeID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tripID))

	s.mock.ExpectExec("^UPDATE \"items\"*").
		WithArgs(AnyTime{}, tripID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	s.mock.ExpectExec("^UPDATE \"grocery_trips\" (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))

	storeUsersID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(storeUsersID))

	s.mock.ExpectExec("UPDATE \"store_users\"*").
		WithArgs(AnyTime{}, storeID, storeUsersID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	store, err := DeleteStore(storeID, userID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), store.ID, storeID)
}
