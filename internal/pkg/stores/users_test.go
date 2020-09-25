package stores

import (
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *Suite) TestInviteToStoreByEmail_UserExistsNotYetAdded() {
	storeID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(storeID))

	email := "test@example.com"
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
	assert.Equal(s.T(), storeUser.Email, email)
}

func (s *Suite) TestInviteToStoreByEmail_UserExistsAlreadyAdded() {
	storeID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(storeID))

	email := "test@example.com"
	storeUserID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID, email, false).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(storeUserID))

	storeUser, err := InviteToStoreByEmail(storeID, email)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), storeUser.ID, storeUserID)
	assert.Equal(s.T(), storeUser.Email, email)
}

func (s *Suite) TestRemoveUserFromStore_StoreNotFound() {
	storeID := uuid.NewV4()
	user := models.User{ID: uuid.NewV4()}
	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, e := RemoveUserFromStore(user, storeID)
	require.Error(s.T(), e)
	assert.Equal(s.T(), e.Error(), "store not found")
}

func (s *Suite) TestRemoveUserFromStore_StoreUserNotFound() {
	storeID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(storeID))

	user := models.User{ID: uuid.NewV4()}
	_, e := RemoveUserFromStore(user, storeID)
	require.Error(s.T(), e)
	assert.Equal(s.T(), e.Error(), "store user not found")
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
	assert.Equal(s.T(), lu.(*models.StoreUser).ID, storeUser.ID)
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
	assert.Equal(s.T(), lu.(*models.StoreUser).ID, storeUser.ID)
}

func (s *Suite) TestAddUserToStore_UserNotFoundInStore() {
	storeID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(storeID))

	user := models.User{ID: uuid.NewV4(), Email: "test@example.com"}
	storeUserID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID, user.Email).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(storeUserID))

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

	s.mock.ExpectBegin()
	s.mock.ExpectExec("^UPDATE \"store_users\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	storeUser, err := AddUserToStore(user, storeID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), storeUser.(*models.StoreUser).ID, storeUserID)
	assert.Equal(s.T(), storeUser.(*models.StoreUser).UserID, user.ID)
	assert.Equal(s.T(), storeUser.(*models.StoreUser).Email, "")
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
		WithArgs(storeID).
		WillReturnRows(storeUserRows)

	storeUsers, err := RetrieveStoreUsers(store.ID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), len(storeUsers.([]models.StoreUser)), 2)
	assert.Equal(s.T(), storeUsers.([]models.StoreUser)[0].StoreID, store.ID)
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
	assert.Equal(s.T(), e.Error(), "record not found")
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
	assert.Equal(s.T(), e.Error(), "record not found")
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
	assert.Equal(s.T(), creatorUser.(*models.User).Email, email)
}
