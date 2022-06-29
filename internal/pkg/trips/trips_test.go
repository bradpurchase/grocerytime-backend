package trips

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type Suite struct {
	suite.Suite

	DB   *gorm.DB
	mock sqlmock.Sqlmock
}

func (s *Suite) SetupSuite() {
	var (
		dbMock *sql.DB
		err    error
	)

	dbMock, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)
	s.DB, err = gorm.Open(postgres.New(postgres.Config{Conn: dbMock}), &gorm.Config{})
	require.NoError(s.T(), err)

	db.Manager = s.DB
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

// Trips

func (s *Suite) TestRetrieveCurrentStoreTripForUser_UserNotMemberOfStore() {
	storeID := uuid.NewV4()
	user := models.User{ID: uuid.NewV4()}
	_, err := RetrieveCurrentStoreTripForUser(storeID, user)
	require.Error(s.T(), err)
	assert.Equal(s.T(), "user is not a member of this store", err.Error())
}

func (s *Suite) TestRetrieveCurrentStoreTripForUser_TripNotAssociatedWithStore() {
	storeID := uuid.NewV4()
	user := models.User{ID: uuid.NewV4(), Email: "test@example.com"}
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID, user.ID, user.Email).
		WillReturnRows(s.mock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))

	_, err := RetrieveCurrentStoreTripForUser(storeID, user)
	require.Error(s.T(), err)
	assert.Equal(s.T(), "could not find trip associated with this store", err.Error())
}

func (s *Suite) TestRetrieveCurrentStoreTripForUser_FoundResult() {
	storeID := uuid.NewV4()
	user := models.User{ID: uuid.NewV4(), Email: "test@example.com"}
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID, user.ID, user.Email).
		WillReturnRows(s.mock.NewRows([]string{"id", "user_id"}).AddRow(uuid.NewV4(), user.ID))

	tripID := uuid.NewV4()
	tripName := "Week of May 31"
	s.mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(storeID, false).
		WillReturnRows(s.mock.NewRows([]string{"id", "name"}).AddRow(tripID, tripName))

	trip, err := RetrieveCurrentStoreTripForUser(storeID, user)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), tripName, trip.Name)
}

func (s *Suite) TestRetrieveTrips_UserNotActive() {
	storeID := uuid.NewV4()
	userID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT count*").
		WithArgs(storeID, userID, true).
		WillReturnRows(s.mock.NewRows([]string{"count"}).AddRow(0))

	_, e := RetrieveTrips(storeID, userID, false)
	require.Error(s.T(), e)
	assert.Equal(s.T(), "user is not active in this store", e.Error())
}

func (s *Suite) TestRetrieveTrips_Found() {
	storeID := uuid.NewV4()
	userID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT count*").
		WithArgs(storeID, userID, true).
		WillReturnRows(s.mock.NewRows([]string{"count"}).AddRow(1))

	rows := s.mock.NewRows([]string{"id"}).
		AddRow(uuid.NewV4()).
		AddRow(uuid.NewV4())
	s.mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(storeID, false).
		WillReturnRows(rows)

	trips, err := RetrieveTrips(storeID, userID, false)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), 2, len(trips))
}

func (s *Suite) TestRetrieveTrip_NotFound() {
	tripID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(tripID).
		WillReturnRows(s.mock.NewRows([]string{}))

	_, e := RetrieveTrip(tripID)
	require.Error(s.T(), e)
	assert.Equal(s.T(), "trip not found", e.Error())
}

func (s *Suite) TestRetrieveTrip_Found() {
	tripID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(tripID).
		WillReturnRows(s.mock.NewRows([]string{"id"}).AddRow(tripID))

	trip, err := RetrieveTrip(tripID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), tripID, trip.ID)
}

func (s *Suite) TestUpdateTrip_TripNotFound() {
	tripID := uuid.NewV4()

	s.mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(tripID).
		WillReturnRows(s.mock.NewRows([]string{}))

	args := map[string]interface{}{"tripId": tripID}
	_, e := UpdateTrip(args)
	require.Error(s.T(), e)
}

// func (s *Suite) TestUpdateTrip_NameUpdate() {
// 	tripID := uuid.NewV4()
// 	args := map[string]interface{}{
// 		"tripId": tripID,
// 		"name":   "My Second Trip",
// 	}
// 	s.mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
// 		WithArgs(tripID).
// 		WillReturnRows(s.mock.NewRows([]string{"id", "name"}).AddRow(tripID, "My First Trip"))

// 	s.mock.ExpectExec("^UPDATE \"grocery_trips\" SET (.+)$").
// 		WillReturnResult(sqlmock.NewResult(1, 1))

// 	_, err := UpdateTrip(args)
// 	require.NoError(s.T(), err)
// 	assert.Equal(s.T(), "My Second Trip", trip.(models.GroceryTrip).Name)
// }

func (s *Suite) TestUpdateTrip_DupeTripName() {
	tripID := uuid.NewV4()
	storeID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(tripID).
		WillReturnRows(s.mock.
			NewRows([]string{"id", "store_id", "name"}).
			AddRow(tripID, storeID, "My First Trip"))

	args := map[string]interface{}{
		"tripId":    tripID,
		"completed": true,
	}

	s.mock.ExpectExec("^UPDATE \"grocery_trips\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))

	s.mock.ExpectBegin()
	// Test case where a trip already exists in this store with this name
	// and assert that it affixes a count after the name
	//
	// Note: this covers the case where a user creates multiple trips in the same day
	currentTime := time.Now()
	tripName := currentTime.Format("Jan 2, 2006")
	likeTripName := fmt.Sprintf("%%%s%%", tripName)
	s.mock.ExpectQuery("^SELECT count*").
		WithArgs(likeTripName, storeID).
		WillReturnRows(s.mock.NewRows([]string{"count"}).AddRow(1))
	finalTripName := fmt.Sprintf("%s (%d)", tripName, 2)
	s.mock.ExpectQuery("^INSERT INTO \"grocery_trips\" (.+)$").
		WithArgs(storeID, finalTripName, false, false, AnyTime{}, AnyTime{}, nil).
		WillReturnRows(s.mock.NewRows([]string{"store_id"}).AddRow(storeID))

	s.mock.ExpectExec("^UPDATE \"items\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	// AddStapleItemsToNewTrip
	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeID).
		WillReturnRows(s.mock.NewRows([]string{"id"}).AddRow(storeID))
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_staple_items\"*").
		WithArgs(storeID).
		WillReturnRows(s.mock.NewRows([]string{}))

	trip, err := UpdateTrip(args)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), true, trip.(models.GroceryTrip).Completed)
	assert.Equal(s.T(), false, trip.(models.GroceryTrip).CopyRemainingItems)
}

func (s *Suite) TestUpdateTrip_MarkCompleted() {
	tripID := uuid.NewV4()
	storeID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(tripID).
		WillReturnRows(s.mock.
			NewRows([]string{"id", "store_id", "name"}).
			AddRow(tripID, storeID, "My First Trip"))

	args := map[string]interface{}{
		"tripId":    tripID,
		"completed": true,
	}

	s.mock.ExpectExec("^UPDATE \"grocery_trips\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))

	s.mock.ExpectBegin()
	currentTime := time.Now()
	tripName := currentTime.Format("Jan 2, 2006")
	likeTripName := fmt.Sprintf("%%%s%%", tripName)
	s.mock.ExpectQuery("^SELECT count*").
		WithArgs(likeTripName, storeID).
		WillReturnRows(s.mock.NewRows([]string{"count"}).AddRow(0))
	s.mock.ExpectQuery("^INSERT INTO \"grocery_trips\" (.+)$").
		WithArgs(storeID, tripName, false, false, AnyTime{}, AnyTime{}, nil).
		WillReturnRows(s.mock.NewRows([]string{"store_id"}).AddRow(storeID))

	s.mock.ExpectExec("^UPDATE \"items\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	// AddStapleItemsToNewTrip
	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeID).
		WillReturnRows(s.mock.NewRows([]string{"id"}).AddRow(storeID))
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_staple_items\"*").
		WithArgs(storeID).
		WillReturnRows(s.mock.NewRows([]string{}))

	trip, err := UpdateTrip(args)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), true, trip.(models.GroceryTrip).Completed)
	assert.Equal(s.T(), false, trip.(models.GroceryTrip).CopyRemainingItems)
}

func (s *Suite) TestUpdateTrip_MarkCompletedAndCopyRemainingItems() {
	tripID := uuid.NewV4()
	storeID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(tripID).
		WillReturnRows(s.mock.NewRows([]string{"id", "store_id", "name"}).AddRow(tripID, storeID, "Trip 1"))

	s.mock.ExpectExec("^UPDATE \"grocery_trips\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))

	s.mock.ExpectBegin()
	newTripID := uuid.NewV4()
	currentTime := time.Now()
	tripName := currentTime.Format("Jan 2, 2006")
	likeTripName := fmt.Sprintf("%%%s%%", tripName)
	s.mock.ExpectQuery("^SELECT count*").
		WithArgs(likeTripName, storeID).
		WillReturnRows(s.mock.NewRows([]string{"count"}).AddRow(0))
	s.mock.ExpectQuery("^INSERT INTO \"grocery_trips\" (.+)$").
		WithArgs(storeID, tripName, false, false, AnyTime{}, AnyTime{}, nil).
		WillReturnRows(s.mock.NewRows([]string{"id"}).AddRow(newTripID))

	// Test creating a category for each remaining item
	remainingItemID := uuid.NewV4()
	itemCategoryID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(tripID, false).
		WillReturnRows(s.mock.NewRows([]string{"id", "category_id"}).AddRow(remainingItemID, itemCategoryID))
	storeCategoryID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_categories\"*").
		WithArgs(itemCategoryID).
		WillReturnRows(s.mock.NewRows([]string{"id", "name"}).AddRow(storeCategoryID, "Produce"))
	s.mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trip_categories\"*").
		WithArgs(newTripID, storeCategoryID).
		WillReturnRows(s.mock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))

		// AddStapleItemsToNewTrip
	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeID).
		WillReturnRows(s.mock.NewRows([]string{"id"}).AddRow(storeID))
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_staple_items\"*").
		WithArgs(storeID).
		WillReturnRows(s.mock.NewRows([]string{}))

	s.mock.ExpectQuery("^INSERT INTO \"items\" (.+)$").
		WithArgs(newTripID, sqlmock.AnyArg(), sqlmock.AnyArg(), nil, sqlmock.AnyArg(), 1, false, 1, nil, sqlmock.AnyArg(), sqlmock.AnyArg(), AnyTime{}, AnyTime{}, nil).
		WillReturnRows(s.mock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))
	s.mock.ExpectExec("^UPDATE \"items\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))

	args := map[string]interface{}{
		"tripId":             tripID,
		"completed":          true,
		"copyRemainingItems": true,
	}
	s.mock.ExpectCommit()

	trip, err := UpdateTrip(args)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), true, trip.(models.GroceryTrip).Completed)
	assert.Equal(s.T(), true, trip.(models.GroceryTrip).CopyRemainingItems)
}

// Items

func (s *Suite) TestRetrieveItems_NoItems() {
	tripID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(tripID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	items, err := RetrieveItems(tripID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), len(items), 0)
}

func (s *Suite) TestRetrieveItems_HasItems() {
	tripID := uuid.NewV4()
	itemRows := sqlmock.
		NewRows([]string{
			"id",
			"grocery_trip_id",
			"name",
			"quantity",
			"completed",
			"notes",
			"created_at",
			"updated_at",
		}).
		AddRow(uuid.NewV4(), tripID, "Apples", 5, false, nil, time.Now(), time.Now()).
		AddRow(uuid.NewV4(), tripID, "Bananas", 2, false, nil, time.Now(), time.Now())
	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(tripID).
		WillReturnRows(itemRows)

	items, err := RetrieveItems(tripID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), 2, len(items))
	assert.Equal(s.T(), tripID, items[0].GroceryTripID)
	assert.Equal(s.T(), "Apples", items[0].Name)
	assert.Equal(s.T(), "Bananas", items[1].Name)
}

func (s *Suite) TestRetrieveItemsInCategory_NoItems() {
	tripID := uuid.NewV4()
	categoryID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(tripID, categoryID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	items, err := RetrieveItemsInCategory(tripID, categoryID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), 0, len(items))
}

func (s *Suite) TestRetrieveItemsInCategory_HasItems() {
	tripID := uuid.NewV4()
	categoryID := uuid.NewV4()
	itemRows := sqlmock.
		NewRows([]string{
			"id",
			"grocery_trip_id",
			"name",
			"quantity",
			"completed",
			"notes",
			"created_at",
			"updated_at",
		}).
		AddRow(uuid.NewV4(), tripID, "Apples", 5, false, nil, time.Now(), time.Now()).
		AddRow(uuid.NewV4(), tripID, "Bananas", 2, false, nil, time.Now(), time.Now())
	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(tripID, categoryID).
		WillReturnRows(itemRows)

	items, err := RetrieveItemsInCategory(tripID, categoryID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), 2, len(items))
	assert.Equal(s.T(), tripID, items[0].GroceryTripID)
	assert.Equal(s.T(), "Apples", items[0].Name)
	assert.Equal(s.T(), "Bananas", items[1].Name)
}

// Add items

func (s *Suite) TestAddItem_TripDoesntExist() {
	tripID := uuid.NewV4()
	userID := uuid.NewV4()
	args := map[string]interface{}{"tripId": tripID, "name": "Test"}

	s.mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(tripID).
		WillReturnRows(s.mock.NewRows([]string{}))

	_, err := AddItem(userID, args)
	require.Error(s.T(), err)
	assert.Equal(s.T(), "trip does not exist", err.Error())
}

func (s *Suite) TestAddItem_UserDoesntBelongInStore() {
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

	args := map[string]interface{}{"tripId": tripID, "name": "Test"}
	_, err := AddItem(userID, args)
	require.Error(s.T(), err)
	assert.Equal(s.T(), "user does not belong to this store", err.Error())
}

func (s *Suite) TestAddItem_NoQuantityArg() {
	userID := uuid.NewV4()
	tripID := uuid.NewV4()
	storeID := uuid.NewV4()
	trip := &models.GroceryTrip{ID: tripID, StoreID: storeID}
	itemName := "Kleenex"
	args := map[string]interface{}{
		"tripId": tripID,
		"name":   itemName,
	}

	s.mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(tripID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "store_id"}).AddRow(trip.ID, trip.StoreID))
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(trip.StoreID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "store_id", "user_id"}).AddRow(uuid.NewV4(), trip.StoreID, userID))

	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_item_category_settings\"*").
		WithArgs(trip.StoreID, strings.ToLower(itemName)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	itemID := uuid.NewV4()
	s.mock.ExpectQuery("^INSERT INTO \"items\" (.+)$").
		WithArgs(trip.ID, sqlmock.AnyArg(), userID, nil, itemName, 1, false, 1, nil, nil, nil, AnyTime{}, AnyTime{}, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(itemID))

	item, err := AddItem(userID, args)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), itemID, item.ID)
	assert.Equal(s.T(), itemName, item.Name)
	assert.Equal(s.T(), 1, item.Quantity)
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

	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_item_category_settings\"*").
		WithArgs(trip.StoreID, "apples").
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	itemID := uuid.NewV4()
	s.mock.ExpectQuery("^INSERT INTO \"items\" (.+)$").
		WithArgs(trip.ID, sqlmock.AnyArg(), userID, nil, "Apples", 6, false, 1, nil, nil, nil, AnyTime{}, AnyTime{}, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(itemID))

	item, err := AddItem(userID, args)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), itemID, item.ID)
	assert.Equal(s.T(), "Apples", item.Name)
	assert.Equal(s.T(), 6, item.Quantity)
}

// Add items to store

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
	assert.Equal(s.T(), "could not find current trip in store", err.Error())
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
	assert.Equal(s.T(), storeID, store.ID)
	assert.Equal(s.T(), userID, store.UserID)
	assert.Equal(s.T(), storeName, store.Name)
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
		WithArgs(storeName, sqlmock.AnyArg(), AnyTime{}, AnyTime{}, nil, userID).
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
	assert.Equal(s.T(), storeName, store.Name)
}

// Update item

func (s *Suite) TestUpdateItem_NoUpdates() {
	itemID := uuid.NewV4()
	tripID := uuid.NewV4()
	storeID := uuid.NewV4()
	userID := uuid.NewV4()

	trip := &models.GroceryTrip{ID: tripID, StoreID: storeID}

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
			AddRow(itemID, trip.ID, userID, "Apples", 5, false, nil, time.Now(), time.Now()))

	s.mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(trip.ID).
		WillReturnRows(s.mock.NewRows([]string{"id", "store_id"}).AddRow(trip.ID, trip.StoreID))
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(trip.StoreID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "store_id", "user_id"}).AddRow(uuid.NewV4(), trip.StoreID, userID))
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_item_category_settings\"*").
		WithArgs(trip.StoreID, "apples").
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(itemID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(itemID))

	s.mock.ExpectExec("^UPDATE \"items\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectExec("^UPDATE \"grocery_trips\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))

	args := map[string]interface{}{"itemId": itemID}
	item, err := UpdateItem(args)
	require.NoError(s.T(), err)
	// Assert no changes
	assert.Equal(s.T(), itemID, item.(*models.Item).ID)
	assert.Equal(s.T(), tripID, item.(*models.Item).GroceryTripID)
	assert.Equal(s.T(), userID, item.(*models.Item).UserID)
	assert.Equal(s.T(), "Apples", item.(*models.Item).Name)
	assert.Equal(s.T(), 5, item.(*models.Item).Quantity)
}

func (s *Suite) TestUpdateItem_UpdateSingleColumn() {
	itemID := uuid.NewV4()
	tripID := uuid.NewV4()
	storeID := uuid.NewV4()
	userID := uuid.NewV4()

	trip := &models.GroceryTrip{ID: tripID, StoreID: storeID}

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
			AddRow(itemID, trip.ID, userID, "Apples", 5, false, nil, time.Now(), time.Now()))

	s.mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(trip.ID).
		WillReturnRows(s.mock.NewRows([]string{"id", "store_id"}).AddRow(trip.ID, trip.StoreID))
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(trip.StoreID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "store_id", "user_id"}).AddRow(uuid.NewV4(), trip.StoreID, userID))
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_item_category_settings\"*").
		WithArgs(trip.StoreID, "apples").
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(itemID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(itemID))

	s.mock.ExpectExec("^UPDATE \"items\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectExec("^UPDATE \"grocery_trips\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))

	completed := true
	args := map[string]interface{}{"itemId": itemID, "completed": completed}

	item, err := UpdateItem(args)
	require.NoError(s.T(), err)
	// Assert only completed state changed
	assert.Equal(s.T(), itemID, item.(*models.Item).ID)
	assert.Equal(s.T(), tripID, item.(*models.Item).GroceryTripID)
	assert.Equal(s.T(), userID, item.(*models.Item).UserID)
	assert.Equal(s.T(), "Apples", item.(*models.Item).Name)
	assert.Equal(s.T(), 5, item.(*models.Item).Quantity)
	assert.Equal(s.T(), &completed, item.(*models.Item).Completed)
}

func (s *Suite) TestUpdateItem_UpdateMultiColumn() {
	itemID := uuid.NewV4()
	tripID := uuid.NewV4()
	storeID := uuid.NewV4()
	userID := uuid.NewV4()

	trip := &models.GroceryTrip{ID: tripID, StoreID: storeID}

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
			AddRow(itemID, trip.ID, userID, "Apples", 5, false, nil, time.Now(), time.Now()))

	s.mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(trip.ID).
		WillReturnRows(s.mock.NewRows([]string{"id", "store_id"}).AddRow(trip.ID, trip.StoreID))
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(trip.StoreID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "store_id", "user_id"}).AddRow(uuid.NewV4(), trip.StoreID, userID))
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_item_category_settings\"*").
		WithArgs(trip.StoreID, "bananas").
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(itemID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(itemID))

	s.mock.ExpectExec("^UPDATE \"items\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectExec("^UPDATE \"grocery_trips\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))

	completed := true
	args := map[string]interface{}{
		"itemId":    itemID,
		"quantity":  10,
		"completed": completed,
		"name":      "Bananas",
	}

	item, err := UpdateItem(args)
	require.NoError(s.T(), err)
	// Assert only quantity and completed states changed
	assert.Equal(s.T(), itemID, item.(*models.Item).ID)
	assert.Equal(s.T(), tripID, item.(*models.Item).GroceryTripID)
	assert.Equal(s.T(), userID, item.(*models.Item).UserID)
	assert.Equal(s.T(), "Bananas", item.(*models.Item).Name)
	assert.Equal(s.T(), 10, item.(*models.Item).Quantity)
	assert.Equal(s.T(), &completed, item.(*models.Item).Completed)
}

// Item reordering

func (s *Suite) TestReorderItem_ReorderItemPosition() {
	itemID := uuid.NewV4()
	userID := uuid.NewV4()
	tripID := uuid.NewV4()
	storeID := uuid.NewV4()

	trip := &models.GroceryTrip{ID: tripID, StoreID: storeID}

	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(itemID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "grocery_trip_id", "user_id", "name"}).AddRow(itemID, tripID, userID, "Apples"))

	s.mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(trip.ID).
		WillReturnRows(s.mock.NewRows([]string{"id", "store_id"}).AddRow(trip.ID, trip.StoreID))

	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(trip.StoreID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "store_id", "user_id"}).AddRow(uuid.NewV4(), trip.StoreID, userID))

	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_item_category_settings\"*").
		WithArgs(trip.StoreID, strings.ToLower("apples")).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(itemID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(itemID))

	s.mock.ExpectExec("^UPDATE \"items\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectExec("^UPDATE \"grocery_trips\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))

	s.mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(tripID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tripID))

	trip, err := ReorderItem(itemID, 4)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), tripID, trip.ID)
}

// Mark item as completed

func (s *Suite) TestMarkItemAsCompleted_CouldNotUpdate() {
	userID := uuid.NewV4()

	s.mock.ExpectExec("^UPDATE \"items\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err := MarkItemAsCompleted("", userID)
	require.Error(s.T(), err)
	assert.Equal(s.T(), "could not update items", err.Error())
}

func (s *Suite) TestMarkItemAsCompleted_Updated() {
	userID := uuid.NewV4()
	name := "Apples"

	s.mock.ExpectExec("^UPDATE \"items\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	itemID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(name, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(itemID))

	_, err := MarkItemAsCompleted(name, userID)
	require.NoError(s.T(), err)
}

// Delete item

func (s *Suite) TestDeleteItem_ItemNotFound() {
	itemID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(itemID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, e := DeleteItem(itemID)
	require.Error(s.T(), e)
	assert.Equal(s.T(), "item not found", e.Error())
}

func (s *Suite) TestDeleteItem_SuccessMoreInCategory() {
	itemID := uuid.NewV4()
	categoryID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(itemID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "category_id"}).AddRow(itemID, categoryID))

	s.mock.ExpectBegin()
	s.mock.ExpectExec("^UPDATE \"items\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	s.mock.ExpectExec("^UPDATE \"grocery_trips\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))

	s.mock.ExpectQuery("^SELECT count*").
		WithArgs(categoryID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	item, err := DeleteItem(itemID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), itemID, item.ID)
}

func (s *Suite) TestDeleteItem_SuccessLastInCategory() {
	itemID := uuid.NewV4()
	categoryID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(itemID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "category_id"}).AddRow(itemID, categoryID))

	s.mock.ExpectBegin()
	s.mock.ExpectExec("^UPDATE \"items\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	s.mock.ExpectExec("^UPDATE \"grocery_trips\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))

	s.mock.ExpectQuery("^SELECT count*").
		WithArgs(categoryID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	s.mock.ExpectBegin()
	s.mock.ExpectExec("^UPDATE \"grocery_trip_categories\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	item, err := DeleteItem(itemID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), itemID, item.ID)
}

// Item search

func (s *Suite) TestSearchForItemByName_NoTripHistory() {
	userID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	_, e := SearchForItemByName("zap", userID)
	require.Error(s.T(), e)
	assert.Equal(s.T(), "no item matches the search term", e.Error())
}

func (s *Suite) TestSearchForItemByName_ItemNotFound() {
	userID := uuid.NewV4()
	tripID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tripID))

	name := "zonk"
	nameArg := fmt.Sprintf("%%%s%%", name)
	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(nameArg, tripID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

	result, err := SearchForItemByName(name, userID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "", result.Name)
}

func (s *Suite) TestSearchForItemByName_Found() {
	userID := uuid.NewV4()
	tripID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tripID))

	name := "app"
	nameArg := fmt.Sprintf("%%%s%%", name)
	rows := sqlmock.
		NewRows([]string{"id", "name"}).
		AddRow(uuid.NewV4(), "Apples")
	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(nameArg, tripID).
		WillReturnRows(rows)

	item, err := SearchForItemByName(name, userID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "Apples", item.Name)
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
