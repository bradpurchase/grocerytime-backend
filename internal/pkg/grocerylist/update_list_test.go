package grocerylist

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateListForUser_NoUpdates(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	listID := uuid.NewV4()
	userID := uuid.NewV4()
	itemRows := sqlmock.
		NewRows([]string{
			"id",
			"user_id",
			"name",
		}).
		AddRow(listID, userID, "My Original List")

	mock.ExpectQuery("^SELECT (.+) FROM \"lists\"*").
		WithArgs(listID, userID).
		WillReturnRows(itemRows)

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE \"lists\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	args := map[string]interface{}{"listId": listID}
	list, err := UpdateListForUser(db, userID, args)
	require.NoError(t, err)
	// Assert no changes
	assert.Equal(t, list.(*models.List).ID, listID)
	assert.Equal(t, list.(*models.List).Name, "My Original List")
}

func TestUpdateListForUser_UpdateSingleColumn(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	listID := uuid.NewV4()
	userID := uuid.NewV4()
	listRows := sqlmock.
		NewRows([]string{
			"id",
			"user_id",
		}).
		AddRow(listID, userID)

	mock.ExpectQuery("^SELECT (.+) FROM \"lists\"*").
		WithArgs(listID, userID).
		WillReturnRows(listRows)

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE \"lists\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	args := map[string]interface{}{"listId": listID, "name": "My Renamed List"}
	list, err := UpdateListForUser(db, userID, args)
	require.NoError(t, err)
	// Assert only completed state changed
	assert.Equal(t, list.(*models.List).ID, listID)
	assert.Equal(t, list.(*models.List).UserID, userID)
	assert.Equal(t, list.(*models.List).Name, "My Renamed List")
}
