package stores

import (
	"database/sql/driver"
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

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestInviteToStoreByEmail_UserExistsNotYetAdded(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: dbMock}), &gorm.Config{})
	require.NoError(t, err)

	storeID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(storeID))

	email := "test@example.com"
	mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID, email, false).
		WillReturnRows(sqlmock.NewRows([]string{}))

	mock.ExpectQuery("^INSERT INTO \"store_users\" (.+)$").
		WithArgs(storeID, sqlmock.AnyArg(), email, false, false, AnyTime{}, AnyTime{}, nil).
		WillReturnRows(sqlmock.NewRows([]string{"store_id"}).AddRow(storeID))
	mock.ExpectQuery("^SELECT \"name\" FROM \"stores\"*").
		WithArgs(storeID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(storeID))

	storeUser, err := InviteToStoreByEmail(db, storeID, email)
	require.NoError(t, err)
	assert.Equal(t, storeUser.Email, email)
}

func TestInviteToStoreByEmail_UserExistsAlreadyAdded(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: dbMock}), &gorm.Config{})
	require.NoError(t, err)

	storeID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(storeID))

	email := "test@example.com"
	storeUserID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID, email, false).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(storeUserID))

	storeUser, err := InviteToStoreByEmail(db, storeID, email)
	require.NoError(t, err)
	assert.Equal(t, storeUser.ID, storeUserID)
	assert.Equal(t, storeUser.Email, email)
}

func TestRemoveUserFromStore_StoreNotFound(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: dbMock}), &gorm.Config{})
	require.NoError(t, err)

	storeID := uuid.NewV4()
	user := models.User{ID: uuid.NewV4()}
	mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, e := RemoveUserFromStore(db, user, storeID)
	require.Error(t, e)
	assert.Equal(t, e.Error(), "store not found")
}

func TestRemoveUserFromStore_StoreUserNotFound(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: dbMock}), &gorm.Config{})
	require.NoError(t, err)

	storeID := uuid.NewV4()
	user := models.User{ID: uuid.NewV4()}
	mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(storeID))

	_, e := RemoveUserFromStore(db, user, storeID)
	require.Error(t, e)
	assert.Equal(t, e.Error(), "store user not found")
}

func TestRemoveUserFromStore_SuccessInvitedUser(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: dbMock}), &gorm.Config{})
	require.NoError(t, err)

	storeID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
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
	mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID, user.ID, user.Email).
		WillReturnRows(rows)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE \"store_users\"*").
		WithArgs(AnyTime{}, storeUser.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Test querying for data to send the email about invite being declined
	creatorUserID := uuid.NewV4()
	creatorUser := models.User{ID: creatorUserID, Email: "creator@example.com"}
	mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID, true).
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(creatorUserID))
	mock.ExpectQuery("^SELECT \"email\" FROM \"users\"*").
		WithArgs(creatorUserID).
		WillReturnRows(sqlmock.NewRows([]string{"email"}).AddRow(creatorUser.Email))

	lu, err := RemoveUserFromStore(db, user, storeID)
	require.NoError(t, err)
	assert.Equal(t, lu.(*models.StoreUser).ID, storeUser.ID)
}

func TestRemoveUserFromStore_SuccessJoinedStoreUser(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: dbMock}), &gorm.Config{})
	require.NoError(t, err)

	storeID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(storeID))

	user := models.User{ID: uuid.NewV4(), Email: "test@example.com"}
	storeUser := models.StoreUser{
		ID:      uuid.NewV4(),
		StoreID: storeID,
		UserID:  user.ID,
	}
	rows := sqlmock.
		NewRows([]string{"id", "user_id"}).
		AddRow(storeUser.ID, storeUser.UserID)
	mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID, user.ID, user.Email).
		WillReturnRows(rows)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE \"store_users\"*").
		WithArgs(AnyTime{}, storeUser.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Test querying for data to send the email about this user leaving the store
	creatorUserID := uuid.NewV4()
	creatorUser := models.User{ID: creatorUserID, Email: "creator@example.com"}
	mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID, true).
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(creatorUserID))
	mock.ExpectQuery("^SELECT \"email\" FROM \"users\"*").
		WithArgs(creatorUserID).
		WillReturnRows(sqlmock.NewRows([]string{"email"}).AddRow(creatorUser.Email))
	mock.ExpectQuery("^SELECT \"email\" FROM \"users\"*").
		WithArgs(storeUser.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"email"}).AddRow(user.Email))

	lu, err := RemoveUserFromStore(db, user, storeID)
	require.NoError(t, err)
	assert.Equal(t, lu.(*models.StoreUser).ID, storeUser.ID)
}

func TestAddUserToStore_UserNotFoundInStore(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: dbMock}), &gorm.Config{})
	require.NoError(t, err)

	storeID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(storeID))

	user := models.User{ID: uuid.NewV4(), Email: "test@example.com"}
	storeUserID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID, user.Email).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(storeUserID))

	_, e := AddUserToStore(db, user, storeID)
	require.Error(t, e)
}

func TestAddUserToStore_Success(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: dbMock}), &gorm.Config{})
	require.NoError(t, err)

	storeID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(storeID))

	email := "test@example.com"
	user := models.User{ID: uuid.NewV4(), Email: email}
	storeUserID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID, user.Email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow(storeUserID, email))

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE \"store_users\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	storeUser, err := AddUserToStore(db, user, storeID)
	require.NoError(t, err)
	assert.Equal(t, storeUser.(*models.StoreUser).ID, storeUserID)
	assert.Equal(t, storeUser.(*models.StoreUser).UserID, user.ID)
	assert.Equal(t, storeUser.(*models.StoreUser).Email, "")
}

func TestRetrieveStoreUsers_HasStoreUsers(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: dbMock}), &gorm.Config{})
	require.NoError(t, err)

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
	mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID).
		WillReturnRows(storeUserRows)

	storeUsers, err := RetrieveStoreUsers(db, store.ID)
	require.NoError(t, err)
	assert.Equal(t, len(storeUsers.([]models.StoreUser)), 2)
	assert.Equal(t, storeUsers.([]models.StoreUser)[0].StoreID, store.ID)
}

func TestRetrieveStoreCreator_StoreUserNotFound(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: dbMock}), &gorm.Config{})
	require.NoError(t, err)

	storeID := uuid.NewV4()
	store := &models.Store{
		ID:        storeID,
		Name:      "Test Store",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	//storeUser := &models.StoreUser{}
	mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID, true).
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, e := RetrieveStoreCreator(db, store.ID)
	require.Error(t, e)
	assert.Equal(t, e.Error(), "record not found")
}

func TestRetrieveStoreCreator_UserNotFound(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: dbMock}), &gorm.Config{})
	require.NoError(t, err)

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
	mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID, true).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "creator"}).AddRow(storeUser.ID, user.ID, storeUser.Creator))
	mock.ExpectQuery("^SELECT (.+) FROM \"users\"*").
		WithArgs(user.ID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, e := RetrieveStoreCreator(db, store.ID)
	require.Error(t, e)
	assert.Equal(t, e.Error(), "record not found")
}

func TestRetrieveStoreCreator_Found(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: dbMock}), &gorm.Config{})
	require.NoError(t, err)

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
	mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID, true).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "creator"}).AddRow(storeUser.ID, user.ID, storeUser.Creator))
	mock.ExpectQuery("^SELECT (.+) FROM \"users\"*").
		WithArgs(user.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow(user.ID, user.Email))

	creatorUser, err := RetrieveStoreCreator(db, store.ID)
	require.NoError(t, err)
	assert.Equal(t, creatorUser.(*models.User).Email, email)
}
