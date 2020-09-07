package stores

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestUpdateStoreForUser_NoUpdates(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: dbMock}), &gorm.Config{})
	require.NoError(t, err)

	storeID := uuid.NewV4()
	userID := uuid.NewV4()
	storeRows := sqlmock.
		NewRows([]string{
			"id",
			"user_id",
			"name",
		}).
		AddRow(storeID, userID, "My Original Store")

	mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeID, userID).
		WillReturnRows(storeRows)

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE \"stores\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectQuery("^SELECT u.email FROM store_users AS su*").
		WithArgs(storeID, false).
		WillReturnRows(sqlmock.NewRows([]string{}))

	args := map[string]interface{}{"storeId": storeID}
	store, err := UpdateStoreForUser(db, userID, args)
	require.NoError(t, err)
	// Assert no changes
	assert.Equal(t, store.(*models.Store).ID, storeID)
	assert.Equal(t, store.(*models.Store).Name, "My Original Store")
}

func TestUpdateStoreForUser_UpdateSingleColumn(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: dbMock}), &gorm.Config{})
	require.NoError(t, err)

	storeID := uuid.NewV4()
	userID := uuid.NewV4()
	storeRows := sqlmock.
		NewRows([]string{
			"id",
			"user_id",
		}).
		AddRow(storeID, userID)

	mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeID, userID).
		WillReturnRows(storeRows)

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE \"stores\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectQuery("^SELECT u.email FROM store_users AS su*").
		WithArgs(storeID, false).
		WillReturnRows(sqlmock.NewRows([]string{}))

	args := map[string]interface{}{"storeId": storeID, "name": "My Renamed Store"}
	store, err := UpdateStoreForUser(db, userID, args)
	require.NoError(t, err)
	// Assert only completed state changed
	assert.Equal(t, store.(*models.Store).ID, storeID)
	assert.Equal(t, store.(*models.Store).UserID, userID)
	assert.Equal(t, store.(*models.Store).Name, "My Renamed Store")
}
