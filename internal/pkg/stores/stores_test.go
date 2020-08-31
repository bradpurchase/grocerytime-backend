package stores

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRetrieveStoreForList_NotFound(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	listID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(listID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, e := RetrieveStoreForList(db, listID)
	require.Error(t, e)
	assert.Equal(t, e.Error(), "record not found")
}

func TestRetrieveStoreForList_Found(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	listID := uuid.NewV4()
	storeName := "LobLaws"
	mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(listID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(uuid.NewV4(), storeName))

	store, err := RetrieveStoreForList(db, listID)
	require.NoError(t, err)
	assert.Equal(t, store.(*models.Store).Name, storeName)
}
