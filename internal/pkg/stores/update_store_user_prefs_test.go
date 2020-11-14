package stores

import (
	"github.com/DATA-DOG/go-sqlmock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *Suite) TestUpdateStoreUserPrefs_UpdateDefaultStore() {
	storeUserID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_user_preferences\"*").
		WithArgs(storeUserID).
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

func (s *Suite) TestUpdateStoreUserPrefs_UpdateMultiColumns() {
	storeUserID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_user_preferences\"*").
		WithArgs(storeUserID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "default_store"}).AddRow(uuid.NewV4(), false))

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
