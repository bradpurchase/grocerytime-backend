package stores

import (
	"github.com/DATA-DOG/go-sqlmock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	assert.Equal(s.T(), storeUserPrefs.DefaultStore, false)
}
