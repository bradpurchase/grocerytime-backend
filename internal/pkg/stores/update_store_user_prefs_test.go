package stores

import (
	"github.com/DATA-DOG/go-sqlmock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	assert.Equal(s.T(), storeUserPrefs.DefaultStore, true)
	assert.Equal(s.T(), storeUserPrefs.Notifications, false)
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
	assert.Equal(s.T(), storeUserPrefs.DefaultStore, true)
	assert.Equal(s.T(), storeUserPrefs.Notifications, false)
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
	assert.Equal(s.T(), storeUserPrefs.DefaultStore, true)
	assert.Equal(s.T(), storeUserPrefs.Notifications, true)
}
