package trips

import (
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *Suite) TestAddItemsToStore_CannotFindStore() {
	userID := uuid.NewV4()
	storeName := "Hanks"
	items := []string{"Apples", "Oranges", "Pears"}
	args := map[string]interface{}{"storeName": storeName, "items": items}

	_, err := AddItemsToStore(userID, args)
	require.Error(s.T(), err)
	assert.Equal(s.T(), err.Error(), "could not find or create store")
}

func (s *Suite) TestAddItemsToStore_CannotFindCurrentTrip() {
	userID := uuid.NewV4()
	storeName := "Hanks"
	storeID := uuid.NewV4()
	items := []string{"Apples", "Oranges", "Pears"}
	args := map[string]interface{}{"storeName": storeName, "items": items}
	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(userID, storeName).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name"}).AddRow(storeID, userID, storeName))

	_, err := AddItemsToStore(userID, args)
	require.Error(s.T(), err)
	assert.Equal(s.T(), err.Error(), "could not find current trip in store")
}

// func (s *Suite) TestAddItemsToStore_CannotCreateItems() {
// 	userID := uuid.NewV4()
// 	storeName := "Hanks"
// 	storeID := uuid.NewV4()
// 	items := []interface{}{"Apples", "Oranges", "Pears"}
// 	args := map[string]interface{}{"storeName": storeName, "items": items}
// 	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
// 		WithArgs(userID, storeName).
// 		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name"}).AddRow(storeID, userID, storeName))

// 	tripID := uuid.NewV4()
// 	trip := models.GroceryTrip{ID: tripID, StoreID: storeID}
// 	s.mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
// 		WithArgs(storeID, false).
// 		WillReturnRows(sqlmock.NewRows([]string{"id", "store_id"}).AddRow(trip.ID, trip.StoreID))
// 	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
// 		WithArgs(trip.StoreID, userID).
// 		WillReturnRows(sqlmock.NewRows([]string{"id", "store_id", "user_id"}))

// 	_, err := AddItemsToStore(userID, args)
// 	require.Error(s.T(), err)
// 	assert.Equal(s.T(), err.Error(), "could not create item")
// }

// func (s *Suite) TestAddItemsToStore_CreatesItems() {
// 	userID := uuid.NewV4()
// 	storeName := "Hanks"
// 	storeID := uuid.NewV4()
// 	items := []interface{}{"Apples", "Oranges", "Pears"}
// 	args := map[string]interface{}{"storeName": storeName, "items": items}
// 	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
// 		WithArgs(userID, storeName).
// 		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name"}).AddRow(storeID, userID, storeName))

// 	tripID := uuid.NewV4()
// 	s.mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
// 		WithArgs(storeID, false).
// 		WillReturnRows(sqlmock.NewRows([]string{"id", "completed"}).AddRow(tripID, false))

// 	itemID := uuid.NewV4()
// 	s.mock.ExpectQuery("^INSERT INTO \"items\" (.+)$").
// 		WithArgs(tripID, sqlmock.AnyArg(), userID, "Apples", 1, false, 1, AnyTime{}, AnyTime{}, nil).
// 		WillReturnRows(sqlmock.NewRows([]string{"id", "trip_id"}).AddRow(itemID, tripID))

// 	addedItems, err := AddItemsToStore(userID, args)
// 	require.NoError(s.T(), err)
// 	assert.Equal(s.T(), len(addedItems), 3)
// }

func (s *Suite) TestFindOrCreateStore_ExistingStoreFound() {
	userID := uuid.NewV4()
	storeID := uuid.NewV4()
	storeName := "Hanks"
	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(userID, storeName).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name"}).AddRow(storeID, userID, storeName))

	store, err := FindOrCreateStore(userID, storeName)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), store.ID, storeID)
	assert.Equal(s.T(), store.UserID, userID)
	assert.Equal(s.T(), store.Name, storeName)
}

func (s *Suite) TestFindOrCreateStore_StoreCreated() {
	userID := uuid.NewV4()
	storeName := "Hanks"
	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(userID, storeName).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	storeID := uuid.NewV4()
	s.mock.ExpectBegin()
	s.mock.ExpectQuery("^INSERT INTO \"stores\" (.+)$").
		WithArgs(storeName, AnyTime{}, AnyTime{}, nil, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(storeID, userID))
	s.mock.ExpectQuery("^INSERT INTO \"store_users\" (.+)$").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "", true, true, AnyTime{}, AnyTime{}, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))

	categories := fetchCategories()
	for i := range categories {
		s.mock.ExpectQuery("^INSERT INTO \"store_categories\" (.+)$").
			WithArgs(sqlmock.AnyArg(), categories[i], AnyTime{}, AnyTime{}, nil).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))
	}

	currentTime := time.Now()
	tripName := currentTime.Format("Jan 02, 2006")
	s.mock.ExpectQuery("^SELECT count*").
		WithArgs(tripName, sqlmock.AnyArg()).
		WillReturnRows(s.mock.NewRows([]string{"count"}).AddRow(0))
	s.mock.ExpectQuery("^INSERT INTO \"grocery_trips\" (.+)$").
		WithArgs(sqlmock.AnyArg(), tripName, false, false, AnyTime{}, AnyTime{}, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))
	s.mock.ExpectCommit()

	store, err := FindOrCreateStore(userID, storeName)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), store.UserID, userID)
	assert.Equal(s.T(), store.Name, storeName)
}

func (s *Suite) TestFindCurrentTripIDInStore_NotFound() {
	storeID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(storeID, false).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	_, err := FindCurrentTripIDInStore(storeID)
	require.Error(s.T(), err)
	assert.Equal(s.T(), err.Error(), "record not found")
}

func (s *Suite) TestFindCurrentTripIDInStore_Found() {
	storeID := uuid.NewV4()
	tripID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(storeID, false).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tripID))

	currentTripID, err := FindCurrentTripIDInStore(storeID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), currentTripID, tripID)
}

// TODO: duplicated code with the store model... DRY this up
func fetchCategories() [20]string {
	categories := [20]string{
		"Produce",
		"Bakery",
		"Meat",
		"Seafood",
		"Dairy",
		"Cereal",
		"Baking",
		"Dry Goods",
		"Canned Goods",
		"Frozen Foods",
		"Cleaning",
		"Paper Products",
		"Beverages",
		"Candy & Snacks",
		"Condiments",
		"Personal Care",
		"Baby",
		"Alcohol",
		"Pharmacy",
		"Misc.",
	}
	return categories
}
