package stores

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type Suite struct {
	suite.Suite

	DB   *gorm.DB
	mock sqlmock.Sqlmock
}

func (s *Suite) SetupSuite() {
	var (
		dbMock *sql.DB
		err    error
	)

	dbMock, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)
	s.DB, err = gorm.Open(postgres.New(postgres.Config{Conn: dbMock}), &gorm.Config{})
	require.NoError(s.T(), err)

	db.Manager = s.DB
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

// Create Store

func (s *Suite) TestCreateStore_DupeStore() {
	userID := uuid.NewV4()
	storeName := "My Dupe Store"
	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeName, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name"}).AddRow(uuid.NewV4(), userID, storeName))

	_, e := CreateStore(userID, storeName)
	require.Error(s.T(), e)
	assert.Equal(s.T(), e.Error(), "you already added a store with this name")
}

func (s *Suite) TestCreateStore_Created() {
	userID := uuid.NewV4()
	storeName := "My New Store"
	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeName, userID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	storeID := uuid.NewV4()
	s.mock.ExpectBegin()
	s.mock.ExpectQuery("^INSERT INTO \"stores\" (.+)$").
		WithArgs(storeName, sqlmock.AnyArg(), AnyTime{}, AnyTime{}, nil, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(storeID, userID))

	storeUserID := uuid.NewV4()
	s.mock.ExpectQuery("^INSERT INTO \"store_users\" (.+)$").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "", true, true, AnyTime{}, AnyTime{}, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(storeUserID))
	s.mock.ExpectQuery("^INSERT INTO \"store_user_preferences\" (.+)$").
		WithArgs(storeUserID, false, true, AnyTime{}, AnyTime{}, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))

	categories := fetchCategories()
	for i := range categories {
		s.mock.ExpectQuery("^INSERT INTO \"store_categories\" (.+)$").
			WithArgs(sqlmock.AnyArg(), categories[i], AnyTime{}, AnyTime{}, nil).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))
	}

	currentTime := time.Now()
	tripName := currentTime.Format("Jan 2, 2006")
	likeTripName := fmt.Sprintf("%%%s%%", tripName)
	s.mock.ExpectQuery("^SELECT count*").
		WithArgs(likeTripName, sqlmock.AnyArg()).
		WillReturnRows(s.mock.NewRows([]string{"count"}).AddRow(0))
	s.mock.ExpectQuery("^INSERT INTO \"grocery_trips\" (.+)$").
		WithArgs(sqlmock.AnyArg(), tripName, false, false, AnyTime{}, AnyTime{}, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))
	s.mock.ExpectCommit()

	store, err := CreateStore(userID, storeName)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), storeName, store.Name)
}

// Update Store

func (s *Suite) TestUpdateStoreForUser_NoUpdates() {
	storeID := uuid.NewV4()
	userID := uuid.NewV4()
	storeRows := sqlmock.
		NewRows([]string{
			"id",
			"user_id",
			"name",
		}).
		AddRow(storeID, userID, "My Original Store")

	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeID, userID).
		WillReturnRows(storeRows)

	s.mock.ExpectBegin()
	s.mock.ExpectExec("^UPDATE \"stores\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	s.mock.ExpectQuery("^SELECT u.email FROM store_users AS su*").
		WithArgs(storeID, false).
		WillReturnRows(sqlmock.NewRows([]string{}))

	args := map[string]interface{}{"storeId": storeID}
	store, err := UpdateStoreForUser(userID, args)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), storeID, store.(*models.Store).ID)
	assert.Equal(s.T(), "My Original Store", store.(*models.Store).Name)
}

func (s *Suite) TestUpdateStoreForUser_UpdateSingleColumn() {
	storeID := uuid.NewV4()
	userID := uuid.NewV4()
	storeRows := sqlmock.
		NewRows([]string{
			"id",
			"user_id",
		}).
		AddRow(storeID, userID)

	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeID, userID).
		WillReturnRows(storeRows)

	s.mock.ExpectBegin()
	s.mock.ExpectExec("^UPDATE \"stores\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	s.mock.ExpectQuery("^SELECT u.email FROM store_users AS su*").
		WithArgs(storeID, false).
		WillReturnRows(sqlmock.NewRows([]string{}))

	args := map[string]interface{}{"storeId": storeID, "name": "My Renamed Store"}
	store, err := UpdateStoreForUser(userID, args)
	require.NoError(s.T(), err)
	// Assert only completed state changed
	assert.Equal(s.T(), storeID, store.(*models.Store).ID, storeID)
	assert.Equal(s.T(), userID, store.(*models.Store).UserID, userID)
	assert.Equal(s.T(), "My Renamed Store", store.(*models.Store).Name)
}

// User stores

func (s *Suite) TestRetrieveUserStores_NoStores() {
	user := models.User{ID: uuid.NewV4(), Email: "test@example.com"}
	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(user.ID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	userStores, err := RetrieveUserStores(user)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), 0, len(userStores))
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
	assert.Equal(s.T(), "Provigo Grocery Store", userStores[0].Name)
	assert.Equal(s.T(), user.ID, userStores[0].UserID)
	assert.Equal(s.T(), "Beer Store Store", userStores[1].Name)
	assert.Equal(s.T(), user.ID, userStores[1].UserID)
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
	assert.Equal(s.T(), "Provigo Grocery Store", userStores[0].Name)
	assert.Equal(s.T(), user.ID, userStores[0].UserID)
	assert.Equal(s.T(), "Beer Store Store", userStores[1].Name)
	assert.Equal(s.T(), "Shared Store", userStores[2].Name)
	assert.Equal(s.T(), sharingUserID, userStores[2].UserID)
}

func (s *Suite) TestRetrieveInvitedUserStores_NoneFound() {
	user := models.User{ID: uuid.NewV4(), Email: "test@example.com"}
	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(user.Email, false).
		WillReturnRows(sqlmock.NewRows([]string{}))

	userStores, err := RetrieveInvitedUserStores(user)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), 0, len(userStores))
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
	assert.Equal(s.T(), 2, len(userStores))
	assert.Equal(s.T(), "Provigo Grocery Store", userStores[0].Name)
	assert.Equal(s.T(), user.ID, userStores[0].UserID)
	assert.Equal(s.T(), "Beer Store Store", userStores[1].Name)
	assert.Equal(s.T(), user.ID, userStores[1].UserID)
}

func (s *Suite) TestRetrieveStoreForUser_NotFound() {
	storeID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, e := RetrieveStoreForUser(storeID, uuid.NewV4())
	require.Error(s.T(), e)
	assert.Equal(s.T(), "record not found", e.Error())
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
	assert.Equal(s.T(), "Example Store", store.Name)
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
	assert.Equal(s.T(), "Example Store", store.Name)
}

func (s *Suite) TestRetrieveStoreForUserByName_NotFound() {
	storeName := "Example Store"
	userID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeName, userID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, e := RetrieveStoreForUserByName(storeName, userID)
	require.Error(s.T(), e)
	assert.Equal(s.T(), "record not found", e.Error())
}

func (s *Suite) TestRetrieveStoreForUserByName_Found() {
	storeName := "Example Store"
	userID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeName, userID).
		WillReturnRows(sqlmock.NewRows([]string{"name", "user_id"}).AddRow(storeName, userID))

	store, err := RetrieveStoreForUserByName(storeName, userID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), storeName, store.Name)
}

func (s *Suite) TestDeleteStore_StoreNotFound() {
	storeID := uuid.NewV4()
	userID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeID, userID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, e := DeleteStore(storeID, userID)
	require.Error(s.T(), e)
	assert.Equal(s.T(), "couldn't retrieve store", e.Error())
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

	tripID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(storeID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tripID))

	s.mock.ExpectExec("^UPDATE \"items\"*").
		WithArgs(AnyTime{}, tripID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	s.mock.ExpectExec("^UPDATE \"grocery_trips\" (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))

	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID, true).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(uuid.NewV4(), userID))

	s.mock.ExpectQuery("^SELECT \"email\" FROM \"users\"*").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))

	s.mock.ExpectExec("^UPDATE \"store_users\" (.+)$").
		WithArgs(AnyTime{}, storeID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	s.mock.ExpectCommit()

	store, err := DeleteStore(storeID, userID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), storeID, store.ID)
}

// Store users

func (s *Suite) TestInviteToStoreByEmail_UserExistsNotYetAdded() {
	storeID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(storeID))

	email := "test@example.com"
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID, email, email).
		WillReturnRows(s.mock.NewRows([]string{"count"}).AddRow(0))
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID, email, false).
		WillReturnRows(sqlmock.NewRows([]string{}))

	s.mock.ExpectQuery("^INSERT INTO \"store_users\" (.+)$").
		WithArgs(storeID, sqlmock.AnyArg(), email, false, false, AnyTime{}, AnyTime{}, nil).
		WillReturnRows(sqlmock.NewRows([]string{"store_id"}).AddRow(storeID))
	s.mock.ExpectQuery("^SELECT name, user_id FROM \"stores\"*").
		WithArgs(storeID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(storeID))
	s.mock.ExpectQuery("^SELECT \"name\" FROM \"users\"*").
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))

	storeUser, err := InviteToStoreByEmail(storeID, email)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), email, storeUser.Email)
}

func (s *Suite) TestInviteToStoreByEmail_UserExistsAlreadyAdded() {
	storeID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(storeID))

	email := "test@example.com"
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID, email, email).
		WillReturnRows(s.mock.NewRows([]string{"count"}).AddRow(1))

	_, err := InviteToStoreByEmail(storeID, email)
	require.Error(s.T(), err)
	assert.Equal(s.T(), "this store is already being shared with this user", err.Error())
}

func (s *Suite) TestRemoveUserFromStore_StoreNotFound() {
	storeID := uuid.NewV4()
	user := models.User{ID: uuid.NewV4()}
	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, e := RemoveUserFromStore(user, storeID)
	require.Error(s.T(), e)
	assert.Equal(s.T(), "store not found", e.Error())
}

func (s *Suite) TestRemoveUserFromStore_StoreUserNotFound() {
	storeID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(storeID))

	user := models.User{ID: uuid.NewV4()}
	_, e := RemoveUserFromStore(user, storeID)
	require.Error(s.T(), e)
	assert.Equal(s.T(), "store user not found", e.Error())
}

func (s *Suite) TestRemoveUserFromStore_SuccessInvitedUser() {
	storeID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(storeID))

	user := models.User{ID: uuid.NewV4(), Email: "test@example.com"}
	storeUser := models.StoreUser{
		ID:      uuid.NewV4(),
		StoreID: storeID,
		Email:   user.Email,
	}
	rows := sqlmock.
		NewRows([]string{"id", "email"}).
		AddRow(storeUser.ID, storeUser.Email)
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID, user.ID, user.Email).
		WillReturnRows(rows)

	s.mock.ExpectBegin()
	s.mock.ExpectExec("UPDATE \"store_users\"*").
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	// Test querying for data to send the email about invite being declined
	creatorUserID := uuid.NewV4()
	creatorUser := models.User{ID: creatorUserID, Email: "creator@example.com"}
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID, true).
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(creatorUserID))
	s.mock.ExpectQuery("^SELECT \"email\" FROM \"users\"*").
		WithArgs(creatorUserID).
		WillReturnRows(sqlmock.NewRows([]string{"email"}).AddRow(creatorUser.Email))

	lu, err := RemoveUserFromStore(user, storeID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), storeUser.ID, lu.(*models.StoreUser).ID)
}

func (s *Suite) TestRemoveUserFromStore_SuccessJoinedStoreUser() {
	storeID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(storeID))

	user := models.User{ID: uuid.NewV4(), Email: "test@example.com", Name: "Jane"}
	storeUser := models.StoreUser{
		ID:      uuid.NewV4(),
		StoreID: storeID,
		UserID:  user.ID,
	}
	rows := sqlmock.
		NewRows([]string{"id", "user_id"}).
		AddRow(storeUser.ID, storeUser.UserID)
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID, user.ID, user.Email).
		WillReturnRows(rows)

	s.mock.ExpectBegin()
	s.mock.ExpectExec("UPDATE \"store_users\"*").
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	// Test querying for data to send the email about this user leaving the store
	creatorUserID := uuid.NewV4()
	creatorUser := models.User{ID: creatorUserID, Email: "creator@example.com"}
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID, true).
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(creatorUserID))
	s.mock.ExpectQuery("^SELECT \"email\" FROM \"users\"*").
		WithArgs(creatorUserID).
		WillReturnRows(sqlmock.NewRows([]string{"email"}).AddRow(creatorUser.Email))
	s.mock.ExpectQuery("^SELECT \"name\" FROM \"users\"*").
		WithArgs(storeUser.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(user.Name))

	lu, err := RemoveUserFromStore(user, storeID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), storeUser.ID, lu.(*models.StoreUser).ID)
}

func (s *Suite) TestAddUserToStoreWithCode_CodeInvalid() {
	user := models.User{ID: uuid.NewV4(), Email: "test@example.com"}

	code := "DEF456"
	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(code).
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, e := AddUserToStoreWithCode(user, code, "Test")
	require.Error(s.T(), e)
	assert.Equal(s.T(), "sorry, that code was invalid", e.Error())
}

func (s *Suite) TestAddUserToStoreWithCode_UserAlreadyInStore() {
	user := models.User{ID: uuid.NewV4(), Email: "test@example.com"}
	storeID := uuid.NewV4()
	storeUser := models.StoreUser{
		ID:      uuid.NewV4(),
		StoreID: storeID,
		UserID:  user.ID,
	}

	code := "ABC123"
	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(code).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(storeID))

	// Assert that record was retrieved, not created (since we call FirstOrCreate)
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID, user.ID).
		WillReturnRows(s.mock.NewRows([]string{"id"}).AddRow(storeUser.ID))

	su, err := AddUserToStoreWithCode(user, code, "Test")
	require.NoError(s.T(), err)
	assert.Equal(s.T(), storeUser.ID, su.ID)
}

func (s *Suite) TestAddUserToStoreWithCode_UserNotAlreadyInStore() {
	user := models.User{ID: uuid.NewV4(), Email: "test@example.com"}
	storeID := uuid.NewV4()

	code := "ABC123"
	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(code).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(storeID))

	// Assert that record was created, not retrieved (since we call FirstOrCreate)
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID, user.ID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	storeUserID := uuid.NewV4()
	s.mock.ExpectQuery("^INSERT INTO \"store_users\" (.+)$").
		WithArgs(storeID, user.ID, "", false, true, AnyTime{}, AnyTime{}, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(storeUserID))
	s.mock.ExpectQuery("^INSERT INTO \"store_user_preferences\" (.+)$").
		WithArgs(storeUserID, false, true, AnyTime{}, AnyTime{}, nil).
		WillReturnRows(sqlmock.NewRows([]string{"store_user_id"}).AddRow(storeUserID))

	su, err := AddUserToStoreWithCode(user, code, "Test")
	require.NoError(s.T(), err)
	assert.Equal(s.T(), storeID, su.StoreID)
	assert.Equal(s.T(), user.ID, su.UserID)
}

// AddUserToStore is DEPRECATED

func (s *Suite) TestAddUserToStore_UserNotFoundInStore() {
	storeID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(storeID))

	user := models.User{ID: uuid.NewV4(), Email: "test@example.com"}
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID, user.Email).
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, e := AddUserToStore(user, storeID)
	require.Error(s.T(), e)
}

func (s *Suite) TestAddUserToStore_Success() {
	storeID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(storeID))

	email := "test@example.com"
	user := models.User{ID: uuid.NewV4(), Email: email}
	storeUserID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID, user.Email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow(storeUserID, email))

	s.mock.ExpectExec("^UPDATE \"store_users\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))

	s.mock.ExpectQuery("^INSERT INTO \"store_user_preferences\" (.+)$").
		WithArgs(storeUserID, false, true, AnyTime{}, AnyTime{}, nil).
		WillReturnRows(sqlmock.NewRows([]string{"store_user_id"}).AddRow(storeUserID))

	storeUser, err := AddUserToStore(user, storeID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), storeUserID, storeUser.ID)
	assert.Equal(s.T(), user.ID, storeUser.UserID)
	assert.Equal(s.T(), "", storeUser.Email)
	assert.Equal(s.T(), false, storeUser.Preferences.DefaultStore)
	assert.Equal(s.T(), false, storeUser.Preferences.Notifications)
}

func (s *Suite) TestRetrieveStoreUsers_HasStoreUsers() {
	storeID := uuid.NewV4()
	store := &models.Store{
		ID:        storeID,
		Name:      "Test Store",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	storeUserRows := sqlmock.
		NewRows([]string{
			"id",
			"store_id",
			"user_id",
			"creator",
			"active",
			"created_at",
			"updated_at",
		}).
		AddRow(uuid.NewV4(), storeID, uuid.NewV4(), true, true, time.Now(), time.Now()).
		AddRow(uuid.NewV4(), storeID, uuid.NewV4(), false, true, time.Now(), time.Now())
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID, true).
		WillReturnRows(storeUserRows)

	storeUsers, err := RetrieveStoreUsers(store.ID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), 2, len(storeUsers))
	assert.Equal(s.T(), store.ID, storeUsers[0].StoreID)
}

func (s *Suite) TestRetrieveStoreCreator_StoreUserNotFound() {
	storeID := uuid.NewV4()
	store := &models.Store{
		ID:        storeID,
		Name:      "Test Store",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	//storeUser := &models.StoreUser{}
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID, true).
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, e := RetrieveStoreCreator(store.ID)
	require.Error(s.T(), e)
	assert.Equal(s.T(), "record not found", e.Error())
}

func (s *Suite) TestRetrieveStoreCreator_UserNotFound() {
	storeID := uuid.NewV4()
	store := &models.Store{
		ID:        storeID,
		Name:      "Test Store",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	userID := uuid.NewV4()
	email := "test@example.com"
	user := &models.User{ID: userID, Email: email}
	storeUserCreator := true

	storeUser := &models.StoreUser{
		ID:      uuid.NewV4(),
		StoreID: storeID,
		UserID:  user.ID,
		Creator: &storeUserCreator,
	}
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID, true).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "creator"}).AddRow(storeUser.ID, user.ID, storeUser.Creator))
	s.mock.ExpectQuery("^SELECT (.+) FROM \"users\"*").
		WithArgs(user.ID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, e := RetrieveStoreCreator(store.ID)
	require.Error(s.T(), e)
	assert.Equal(s.T(), "record not found", e.Error())
}

func (s *Suite) TestRetrieveStoreCreator_Found() {
	storeID := uuid.NewV4()
	store := &models.Store{
		ID:        storeID,
		Name:      "Test Store",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	userID := uuid.NewV4()
	email := "test@example.com"
	user := &models.User{ID: userID, Email: email}
	storeUserCreator := true

	storeUser := &models.StoreUser{
		ID:      uuid.NewV4(),
		StoreID: storeID,
		UserID:  user.ID,
		Creator: &storeUserCreator,
	}
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID, true).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "creator"}).AddRow(storeUser.ID, user.ID, storeUser.Creator))
	s.mock.ExpectQuery("^SELECT (.+) FROM \"users\"*").
		WithArgs(user.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow(user.ID, user.Email))

	creatorUser, err := RetrieveStoreCreator(store.ID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), email, creatorUser.(*models.User).Email)
}

// TestUpdateStoreUserPrefs_UpdateDefaultStoreMultiStores tests the case where
// a store is updated to become the default for a StoreUser record,
// and the user's other stores are unmarked as default
func (s *Suite) TestUpdateStoreUserPrefs_UpdateDefaultStoreMultiStores() {
	storeUserID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_user_preferences\"*").
		WithArgs(storeUserID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))

	s.mock.ExpectExec("^UPDATE \"store_user_preferences\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// AfterUpdate hook to unmark defaults for other store users
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))
	s.mock.ExpectExec("^UPDATE \"store_user_preferences\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))

	storeID := uuid.NewV4()
	args := map[string]interface{}{
		"storeId":      storeID,
		"defaultStore": true,
	}
	storeUserPrefs, err := UpdateStoreUserPrefs(storeUserID, args)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), true, storeUserPrefs.DefaultStore)
	assert.Equal(s.T(), false, storeUserPrefs.Notifications)
}

// TestUpdateStoreUserPrefs_UpdateDefaultStoreOnlyStores tests the case where
// a store is updated to become the default for a StoreUser record,
// and the user has no other stores to unmark as default
func (s *Suite) TestUpdateStoreUserPrefs_UpdateDefaultStoreOnlyStores() {
	storeUserID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_user_preferences\"*").
		WithArgs(storeUserID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))

	s.mock.ExpectExec("^UPDATE \"store_user_preferences\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// AfterUpdate hook to unmark defaults for other store users
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{}))

	storeID := uuid.NewV4()
	args := map[string]interface{}{
		"storeId":      storeID,
		"defaultStore": true,
	}
	storeUserPrefs, err := UpdateStoreUserPrefs(storeUserID, args)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), true, storeUserPrefs.DefaultStore)
	assert.Equal(s.T(), false, storeUserPrefs.Notifications)
}

func (s *Suite) TestUpdateStoreUserPrefs_UpdateMultiColumns() {
	storeUserID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_user_preferences\"*").
		WithArgs(storeUserID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "default_store"}).AddRow(uuid.NewV4(), false))

	s.mock.ExpectExec("^UPDATE \"store_user_preferences\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// AfterUpdate hook to unmark defaults for other store users
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))
	s.mock.ExpectExec("^UPDATE \"store_user_preferences\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))

	storeID := uuid.NewV4()
	args := map[string]interface{}{
		"storeId":       storeID,
		"defaultStore":  true,
		"notifications": true,
	}
	storeUserPrefs, err := UpdateStoreUserPrefs(storeUserID, args)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), true, storeUserPrefs.DefaultStore)
	assert.Equal(s.T(), true, storeUserPrefs.Notifications)
}

func (s *Suite) TestRetrieveStoreUserPrefs_NotFound() {
	storeUserID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_user_preferences\"*").
		WithArgs(storeUserID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, e := RetrieveStoreUserPrefs(storeUserID)
	require.Error(s.T(), e)
	assert.Equal(s.T(), e.Error(), "record not found")
}

func (s *Suite) TestRetrieveStoreUserPrefs_Found() {
	storeUserID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_user_preferences\"*").
		WithArgs(storeUserID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))

	storeUserPrefs, err := RetrieveStoreUserPrefs(storeUserID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), false, storeUserPrefs.DefaultStore)
}

// Staple items

func (s *Suite) TestSaveStapleItem_ItemNotFound() {
	itemID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(itemID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	storeID := uuid.NewV4()
	_, err := SaveStapleItem(storeID, itemID)
	require.Error(s.T(), err)
	assert.Equal(s.T(), "record not found", err.Error())
}

func (s *Suite) TestSaveStapleItem_FindExistingStapleItem() {
	itemID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(itemID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(itemID))

	stapleItemID := uuid.NewV4()
	storeID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_staple_items\"*").
		WithArgs(storeID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "store_id"}).AddRow(stapleItemID, storeID))

	s.mock.ExpectExec("^UPDATE \"items\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))

	stapleItem, err := SaveStapleItem(storeID, itemID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), stapleItemID, stapleItem.ID)
}

func (s *Suite) TestSaveStapleItem_CreateNewStapleItem() {
	itemID := uuid.NewV4()
	itemName := "Apples"
	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(itemID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(itemID, itemName))

	stapleItemID := uuid.NewV4()
	storeID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_staple_items\"*").
		WithArgs(storeID, itemName).
		WillReturnRows(sqlmock.NewRows([]string{}))
	s.mock.ExpectQuery("^INSERT INTO \"store_staple_items\" (.+)$").
		WithArgs(storeID, itemName, AnyTime{}, AnyTime{}, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(stapleItemID))

	s.mock.ExpectExec("^UPDATE \"items\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))

	stapleItem, err := SaveStapleItem(storeID, itemID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), stapleItemID, stapleItem.ID)
}

func (s *Suite) TestRemoveStapleItem_StapleItemNotFound() {
	itemID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(itemID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, err := RemoveStapleItem(itemID)
	require.Error(s.T(), err)
	assert.Equal(s.T(), "record not found", err.Error())
}

func (s *Suite) TestRemoveStapleItem_StapleItemRemoved() {
	itemID := uuid.NewV4()
	stapleItemID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(itemID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "staple_item_id"}).AddRow(itemID, stapleItemID))

	s.mock.ExpectExec("^UPDATE \"store_staple_items\"*").
		WillReturnResult(sqlmock.NewResult(1, 1))

	s.mock.ExpectExec("^UPDATE \"items\"*").
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err := RemoveStapleItem(itemID)
	require.NoError(s.T(), err)
}

// TODO: duplicated code with the store model... DRY this up
func fetchCategories() [20]string {
	categories := [20]string{
		"Produce",
		"Bakery",
		"Meat",
		"Seafood",
		"Dairy",
		"Cereal",
		"Baking",
		"Dry Goods",
		"Canned Goods",
		"Frozen Foods",
		"Cleaning",
		"Paper Products",
		"Beverages",
		"Candy & Snacks",
		"Condiments",
		"Personal Care",
		"Baby",
		"Alcohol",
		"Pharmacy",
		"Misc.",
	}
	return categories
}
