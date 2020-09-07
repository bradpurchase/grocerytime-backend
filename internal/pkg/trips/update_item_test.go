package trips

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestUpdateItem_NoUpdates(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	itemID := uuid.NewV4()
	tripID := uuid.NewV4()
	userID := uuid.NewV4()

	mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(itemID).
		WillReturnRows(sqlmock.
			NewRows([]string{
				"id",
				"grocery_trip_id",
				"user_id",
				"name",
				"quantity",
				"completed",
				"created_at",
				"updated_at",
			}).
			AddRow(itemID, tripID, userID, "Apples", 5, false, time.Now(), time.Now()))

	mock.ExpectBegin()
	mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(itemID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(itemID))
	mock.ExpectExec("^UPDATE \"items\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("^UPDATE \"grocery_trips\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	args := map[string]interface{}{"itemId": itemID}
	item, err := UpdateItem(db, args)
	require.NoError(t, err)
	// Assert no changes
	assert.Equal(t, item.(*models.Item).ID, itemID)
	assert.Equal(t, item.(*models.Item).GroceryTripID, tripID)
	assert.Equal(t, item.(*models.Item).UserID, userID)
	assert.Equal(t, item.(*models.Item).Name, "Apples")
	assert.Equal(t, item.(*models.Item).Quantity, 5)
}

func TestUpdateItem_UpdateSingleColumn(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	itemID := uuid.NewV4()
	tripID := uuid.NewV4()
	userID := uuid.NewV4()

	mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(itemID).
		WillReturnRows(sqlmock.
			NewRows([]string{
				"id",
				"grocery_trip_id",
				"user_id",
				"name",
				"quantity",
				"completed",
				"created_at",
				"updated_at",
			}).
			AddRow(itemID, tripID, userID, "Apples", 5, false, time.Now(), time.Now()))

	mock.ExpectBegin()
	mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(itemID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(itemID))
	mock.ExpectExec("^UPDATE \"items\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("^UPDATE \"grocery_trips\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	completed := true
	args := map[string]interface{}{"itemId": itemID, "completed": completed}

	item, err := UpdateItem(db, args)
	require.NoError(t, err)
	// Assert only completed state changed
	assert.Equal(t, item.(*models.Item).ID, itemID)
	assert.Equal(t, item.(*models.Item).GroceryTripID, tripID)
	assert.Equal(t, item.(*models.Item).UserID, userID)
	assert.Equal(t, item.(*models.Item).Name, "Apples")
	assert.Equal(t, item.(*models.Item).Quantity, 5)
	assert.Equal(t, item.(*models.Item).Completed, &completed)
}

func TestUpdateItem_UpdateMultiColumn(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	itemID := uuid.NewV4()
	tripID := uuid.NewV4()
	userID := uuid.NewV4()

	mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(itemID).
		WillReturnRows(sqlmock.
			NewRows([]string{
				"id",
				"grocery_trip_id",
				"user_id",
				"name",
				"quantity",
				"completed",
				"created_at",
				"updated_at",
			}).
			AddRow(itemID, tripID, userID, "Apples", 5, false, time.Now(), time.Now()))

	mock.ExpectBegin()
	mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(itemID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(itemID))
	mock.ExpectExec("^UPDATE \"items\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("^UPDATE \"grocery_trips\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	completed := true
	args := map[string]interface{}{
		"itemId":    itemID,
		"quantity":  10,
		"completed": completed,
		"name":      "Bananas",
	}

	item, err := UpdateItem(db, args)
	require.NoError(t, err)
	// Assert only quantity and completed states changed
	assert.Equal(t, item.(*models.Item).ID, itemID)
	assert.Equal(t, item.(*models.Item).GroceryTripID, tripID)
	assert.Equal(t, item.(*models.Item).UserID, userID)
	assert.Equal(t, item.(*models.Item).Name, "Bananas")
	assert.Equal(t, item.(*models.Item).Quantity, 10)
	assert.Equal(t, item.(*models.Item).Completed, &completed)
}
