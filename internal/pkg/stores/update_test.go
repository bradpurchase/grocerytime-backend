package stores

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	assert.Equal(s.T(), store.(*models.Store).ID, storeID)
	assert.Equal(s.T(), store.(*models.Store).Name, "My Original Store")
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
	assert.Equal(s.T(), store.(*models.Store).ID, storeID)
	assert.Equal(s.T(), store.(*models.Store).UserID, userID)
	assert.Equal(s.T(), store.(*models.Store).Name, "My Renamed Store")
}
