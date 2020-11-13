package trips

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *Suite) TestAddItem_TripDoesntExist() {
	tripID := uuid.NewV4()
	userID := uuid.NewV4()
	args := map[string]interface{}{"tripId": tripID}

	_, err := AddItem(userID, args)
	require.Error(s.T(), err)
	assert.Equal(s.T(), err.Error(), "trip does not exist")
}

func (s *Suite) TestAddItem_UserDoesntBelongInList() {
	userID := uuid.NewV4()
	tripID := uuid.NewV4()
	storeID := uuid.NewV4()
	trip := &models.GroceryTrip{ID: tripID, StoreID: storeID}
	s.mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(tripID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "store_id"}).AddRow(trip.ID, trip.StoreID))
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(trip.StoreID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "store_id", "user_id"}))

	args := map[string]interface{}{"tripId": tripID}
	_, err := AddItem(userID, args)
	require.Error(s.T(), err)
	assert.Equal(s.T(), err.Error(), "user does not belong to this store")
}

func (s *Suite) TestAddItem_AddsItemToTripWithCategoryName() {
	userID := uuid.NewV4()
	tripID := uuid.NewV4()
	storeID := uuid.NewV4()
	trip := &models.GroceryTrip{ID: tripID, StoreID: storeID}
	args := map[string]interface{}{
		"tripId":       tripID,
		"name":         "Apples",
		"quantity":     5,
		"categoryName": "Produce",
	}

	s.mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(tripID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "store_id"}).AddRow(trip.ID, trip.StoreID))
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(trip.StoreID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "store_id", "user_id"}).AddRow(uuid.NewV4(), trip.StoreID, userID))

	itemID := uuid.NewV4()
	s.mock.ExpectQuery("^INSERT INTO \"items\" (.+)$").
		WithArgs(trip.ID, sqlmock.AnyArg(), userID, "Apples", 5, false, 1, AnyTime{}, AnyTime{}, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(itemID))

	item, err := AddItem(userID, args)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), item.ID, itemID)
	assert.Equal(s.T(), item.Name, args["name"])
}

func (s *Suite) TestAddItem_NoQuantityArg() {
	userID := uuid.NewV4()
	tripID := uuid.NewV4()
	storeID := uuid.NewV4()
	trip := &models.GroceryTrip{ID: tripID, StoreID: storeID}
	args := map[string]interface{}{
		"tripId":       tripID,
		"name":         "Kleenex",
		"categoryName": "Cleaning",
	}

	s.mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(tripID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "store_id"}).AddRow(trip.ID, trip.StoreID))
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(trip.StoreID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "store_id", "user_id"}).AddRow(uuid.NewV4(), trip.StoreID, userID))

	itemID := uuid.NewV4()
	s.mock.ExpectQuery("^INSERT INTO \"items\" (.+)$").
		WithArgs(trip.ID, sqlmock.AnyArg(), userID, "Kleenex", 1, false, 1, AnyTime{}, AnyTime{}, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(itemID))

	item, err := AddItem(userID, args)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), item.ID, itemID)
	assert.Equal(s.T(), item.Name, "Kleenex")
	assert.Equal(s.T(), item.Quantity, 1)
}

func (s *Suite) TestAddItem_InlineQuantityInItemName() {
	userID := uuid.NewV4()
	tripID := uuid.NewV4()
	storeID := uuid.NewV4()
	trip := &models.GroceryTrip{ID: tripID, StoreID: storeID}
	args := map[string]interface{}{
		"tripId":       tripID,
		"name":         "Apples x 6",
		"categoryName": "Produce",
	}

	s.mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(tripID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "store_id"}).AddRow(trip.ID, trip.StoreID))
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(trip.StoreID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "store_id", "user_id"}).AddRow(uuid.NewV4(), trip.StoreID, userID))

	itemID := uuid.NewV4()
	s.mock.ExpectQuery("^INSERT INTO \"items\" (.+)$").
		WithArgs(trip.ID, sqlmock.AnyArg(), userID, "Apples", 6, false, 1, AnyTime{}, AnyTime{}, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(itemID))

	item, err := AddItem(userID, args)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), item.ID, itemID)
	assert.Equal(s.T(), item.Name, "Apples")
	assert.Equal(s.T(), item.Quantity, 6)
}
