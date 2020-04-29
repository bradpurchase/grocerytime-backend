package grocerylist

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateItem_NoUpdates(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	itemID := uuid.NewV4()
	listID := uuid.NewV4()
	userID := uuid.NewV4()
	itemRows := sqlmock.
		NewRows([]string{
			"id",
			"list_id",
			"user_id",
			"name",
			"quantity",
			"completed",
			"created_at",
			"updated_at",
		}).
		AddRow(itemID, listID, userID, "Apples", 5, false, time.Now(), time.Now())

	mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(itemID).
		WillReturnRows(itemRows)

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE \"items\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	args := map[string]interface{}{"itemId": itemID}
	item, err := UpdateItem(db, args)
	require.NoError(t, err)
	// Assert no changes
	assert.Equal(t, item.(*models.Item).ID, itemID)
	assert.Equal(t, item.(*models.Item).ListID, listID)
	assert.Equal(t, item.(*models.Item).UserID, userID)
	assert.Equal(t, item.(*models.Item).Name, "Apples")
	assert.Equal(t, item.(*models.Item).Quantity, 5)
	assert.Equal(t, item.(*models.Item).Completed, false)
}

func TestUpdateItem_UpdateSingleColumn(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	itemID := uuid.NewV4()
	listID := uuid.NewV4()
	userID := uuid.NewV4()
	itemRows := sqlmock.
		NewRows([]string{
			"id",
			"list_id",
			"user_id",
			"name",
			"quantity",
			"completed",
			"created_at",
			"updated_at",
		}).
		AddRow(itemID, listID, userID, "Apples", 5, false, time.Now(), time.Now())

	mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(itemID).
		WillReturnRows(itemRows)

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE \"items\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	args := map[string]interface{}{"itemId": itemID, "completed": true}
	item, err := UpdateItem(db, args)
	require.NoError(t, err)
	// Assert only completed state changed
	assert.Equal(t, item.(*models.Item).ID, itemID)
	assert.Equal(t, item.(*models.Item).ListID, listID)
	assert.Equal(t, item.(*models.Item).UserID, userID)
	assert.Equal(t, item.(*models.Item).Name, "Apples")
	assert.Equal(t, item.(*models.Item).Quantity, 5)
	assert.Equal(t, item.(*models.Item).Completed, true)
}

func TestUpdateItem_UpdateMultiColumn(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	itemID := uuid.NewV4()
	listID := uuid.NewV4()
	userID := uuid.NewV4()
	itemRows := sqlmock.
		NewRows([]string{
			"id",
			"list_id",
			"user_id",
			"name",
			"quantity",
			"completed",
			"created_at",
			"updated_at",
		}).
		AddRow(itemID, listID, userID, "Apples", 5, false, time.Now(), time.Now())

	mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(itemID).
		WillReturnRows(itemRows)

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE \"items\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	args := map[string]interface{}{"itemId": itemID, "quantity": 10, "completed": true}
	item, err := UpdateItem(db, args)
	require.NoError(t, err)
	// Assert only quantity and completed states changed
	assert.Equal(t, item.(*models.Item).ID, itemID)
	assert.Equal(t, item.(*models.Item).ListID, listID)
	assert.Equal(t, item.(*models.Item).UserID, userID)
	assert.Equal(t, item.(*models.Item).Name, "Apples")
	assert.Equal(t, item.(*models.Item).Quantity, 10)
	assert.Equal(t, item.(*models.Item).Completed, true)
}
