package stores

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateStore_DupeStore(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	userID := uuid.NewV4()
	storeName := "My Dupe Store"
	mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeName, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name"}).AddRow(uuid.NewV4(), userID, storeName))

	_, e := CreateStore(db, userID, storeName)
	require.Error(t, e)
	assert.Equal(t, e.Error(), "You already added a store with this name")
}

func TestCreateStore_Created(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	userID := uuid.NewV4()
	storeName := "My New Store"
	mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeName, userID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	storeID := uuid.NewV4()
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO \"stores\" (.+)$").
		WithArgs(userID, storeName, AnyTime{}, AnyTime{}, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(storeID))
	mock.ExpectQuery("^INSERT INTO \"store_users\" (.+)$").
		WithArgs(storeID, userID, "", true, true, AnyTime{}, AnyTime{}, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))
	categoryID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"categories\"*").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(categoryID, "Produce"))
	mock.ExpectQuery("^INSERT INTO \"store_categories\" (.+)$").
		WithArgs(storeID, categoryID, AnyTime{}, AnyTime{}, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))

	mock.ExpectQuery("^INSERT INTO \"grocery_trips\" (.+)$").
		WithArgs(storeID, "Trip 1", AnyTime{}, AnyTime{}, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))
	mock.ExpectCommit()

	store, err := CreateStore(db, userID, storeName)
	require.NoError(t, err)
	assert.Equal(t, store.Name, storeName)
}
