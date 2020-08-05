package grocerylist

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateList_DupeList(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	userID := uuid.NewV4()
	listName := "My Dupe List"
	mock.ExpectQuery("^SELECT (.+) FROM \"lists\"*").
		WithArgs(listName, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name"}).AddRow(uuid.NewV4(), userID, listName))

	_, e := CreateList(db, userID, listName)
	require.Error(t, e)
	assert.Equal(t, e.Error(), "You already have a list with this name")
}

func TestCreateList_Created(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	userID := uuid.NewV4()
	listName := "My New List"
	mock.ExpectQuery("^SELECT (.+) FROM \"lists\"*").
		WithArgs(listName, userID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	listID := uuid.NewV4()
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO \"lists\" (.+)$").
		WithArgs(userID, listName, AnyTime{}, AnyTime{}, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(listID))
	mock.ExpectQuery("^INSERT INTO \"list_users\" (.+)$").
		WithArgs(listID, userID, "", true, true, AnyTime{}, AnyTime{}, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))
	mock.ExpectQuery("^INSERT INTO \"grocery_trips\" (.+)$").
		WithArgs(listID, "Trip 1", AnyTime{}, AnyTime{}, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))
	mock.ExpectCommit()

	list, err := CreateList(db, userID, listName)
	require.NoError(t, err)
	assert.Equal(t, list.Name, listName)
}
