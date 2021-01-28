package trips

import (
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *Suite) TestUpdateItem_NoUpdates() {
	itemID := uuid.NewV4()
	tripID := uuid.NewV4()
	userID := uuid.NewV4()

	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(itemID).
		WillReturnRows(sqlmock.
			NewRows([]string{
				"id",
				"grocery_trip_id",
				"user_id",
				"name",
				"quantity",
				"completed",
				"notes",
				"created_at",
				"updated_at",
			}).
			AddRow(itemID, tripID, userID, "Apples", 5, false, nil, time.Now(), time.Now()))

	s.mock.ExpectBegin()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(itemID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(itemID))
	s.mock.ExpectExec("^UPDATE \"items\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectExec("^UPDATE \"grocery_trips\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	args := map[string]interface{}{"itemId": itemID}
	item, err := UpdateItem(args)
	require.NoError(s.T(), err)
	// Assert no changes
	assert.Equal(s.T(), item.(*models.Item).ID, itemID)
	assert.Equal(s.T(), item.(*models.Item).GroceryTripID, tripID)
	assert.Equal(s.T(), item.(*models.Item).UserID, userID)
	assert.Equal(s.T(), item.(*models.Item).Name, "Apples")
	assert.Equal(s.T(), item.(*models.Item).Quantity, 5)
}

// func (s *Suite) TestUpdateItem_UpdateSingleColumn() {
// 	itemID := uuid.NewV4()
// 	tripID := uuid.NewV4()
// 	userID := uuid.NewV4()

// 	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
// 		WithArgs(itemID).
// 		WillReturnRows(sqlmock.
// 			NewRows([]string{
// 				"id",
// 				"grocery_trip_id",
// 				"user_id",
// 				"name",
// 				"quantity",
// 				"completed",
// 				"created_at",
// 				"updated_at",
// 			}).
// 			AddRow(itemID, tripID, userID, "Apples", 5, false, time.Now(), time.Now()))

// 	s.mock.ExpectBegin()
// 	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
// 		WithArgs(itemID).
// 		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(itemID))
// 	s.mock.ExpectExec("^UPDATE \"items\" SET (.+)$").
// 		WillReturnResult(sqlmock.NewResult(1, 1))
// 	s.mock.ExpectExec("^UPDATE \"grocery_trips\" SET (.+)$").
// 		WillReturnResult(sqlmock.NewResult(1, 1))
// 	s.mock.ExpectCommit()

// 	completed := true
// 	args := map[string]interface{}{"itemId": itemID, "completed": completed}

// 	item, err := UpdateItem(args)
// 	require.NoError(s.T(), err)
// 	// Assert only completed state changed
// 	assert.Equal(s.T(), item.(*models.Item).ID, itemID)
// 	assert.Equal(s.T(), item.(*models.Item).GroceryTripID, tripID)
// 	assert.Equal(s.T(), item.(*models.Item).UserID, userID)
// 	assert.Equal(s.T(), item.(*models.Item).Name, "Apples")
// 	assert.Equal(s.T(), item.(*models.Item).Quantity, 5)
// 	assert.Equal(s.T(), item.(*models.Item).Completed, &completed)
// }

// func (s *Suite) TestUpdateItem_UpdateMultiColumn() {
// 	itemID := uuid.NewV4()
// 	tripID := uuid.NewV4()
// 	userID := uuid.NewV4()

// 	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
// 		WithArgs(itemID).
// 		WillReturnRows(sqlmock.
// 			NewRows([]string{
// 				"id",
// 				"grocery_trip_id",
// 				"user_id",
// 				"name",
// 				"quantity",
// 				"completed",
// 				"created_at",
// 				"updated_at",
// 			}).
// 			AddRow(itemID, tripID, userID, "Apples", 5, false, time.Now(), time.Now()))

// 	s.mock.ExpectBegin()
// 	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
// 		WithArgs(itemID).
// 		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(itemID))
// 	s.mock.ExpectExec("^UPDATE \"items\" SET (.+)$").
// 		WillReturnResult(sqlmock.NewResult(1, 1))
// 	s.mock.ExpectExec("^UPDATE \"grocery_trips\" SET (.+)$").
// 		WillReturnResult(sqlmock.NewResult(1, 1))
// 	s.mock.ExpectCommit()

// 	completed := true
// 	args := map[string]interface{}{
// 		"itemId":    itemID,
// 		"quantity":  10,
// 		"completed": completed,
// 		"name":      "Bananas",
// 	}

// 	item, err := UpdateItem(args)
// 	require.NoError(s.T(), err)
// 	// Assert only quantity and completed states changed
// 	assert.Equal(s.T(), item.(*models.Item).ID, itemID)
// 	assert.Equal(s.T(), item.(*models.Item).GroceryTripID, tripID)
// 	assert.Equal(s.T(), item.(*models.Item).UserID, userID)
// 	assert.Equal(s.T(), item.(*models.Item).Name, "Bananas")
// 	assert.Equal(s.T(), item.(*models.Item).Quantity, 10)
// 	assert.Equal(s.T(), item.(*models.Item).Completed, &completed)
// }
