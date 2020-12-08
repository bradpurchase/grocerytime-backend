package trips

import (
	"fmt"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

	storeUserID := uuid.NewV4()
	s.mock.ExpectQuery("^INSERT INTO \"store_users\" (.+)$").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "", true, true, AnyTime{}, AnyTime{}, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(storeUserID))

	s.mock.ExpectQuery("^INSERT INTO \"store_user_preferences\" (.+)$").
		WithArgs(storeUserID, false, true, AnyTime{}, AnyTime{}, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))

	categories := fetchCategories()
	for i := range categories {
		s.mock.ExpectQuery("^INSERT INTO \"store_categories\" (.+)$").
			WithArgs(sqlmock.AnyArg(), categories[i], AnyTime{}, AnyTime{}, nil).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))
	}

	currentTime := time.Now()
	tripName := currentTime.Format("Jan 2, 2006")
	likeTripName := fmt.Sprintf("%%%s%%", tripName)
	s.mock.ExpectQuery("^SELECT count*").
		WithArgs(likeTripName, sqlmock.AnyArg()).
		WillReturnRows(s.mock.NewRows([]string{"count"}).AddRow(0))
	s.mock.ExpectQuery("^INSERT INTO \"grocery_trips\" (.+)$").
		WithArgs(sqlmock.AnyArg(), tripName, false, false, AnyTime{}, AnyTime{}, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))
	s.mock.ExpectCommit()

	store, err := FindOrCreateStore(userID, storeName)
	require.NoError(s.T(), err)
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
