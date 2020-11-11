package trips

import (
	"fmt"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *Suite) TestUpdateTrip_TripNotFound() {
	tripID := uuid.NewV4()

	args := map[string]interface{}{"tripId": tripID}
	result, e := UpdateTrip(args)
	require.Error(s.T(), e)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), e.Error(), "trip does not exist")
}

func (s *Suite) TestUpdateTrip_NameUpdate() {
	tripID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(tripID).
		WillReturnRows(s.mock.NewRows([]string{"id", "name"}).AddRow(tripID, "My First Trip"))

	args := map[string]interface{}{
		"tripId": tripID,
		"name":   "My Second Trip",
	}

	s.mock.ExpectBegin()
	s.mock.ExpectExec("^UPDATE \"grocery_trips\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	trip, err := UpdateTrip(args)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), trip.(models.GroceryTrip).Name, "My Second Trip")
}

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
	s.mock.ExpectQuery("^SELECT count*").
		WithArgs(storeID).
		WillReturnRows(s.mock.NewRows([]string{"count"}).AddRow(0))

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

	trip, err := UpdateTrip(args)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), trip.(models.GroceryTrip).Completed, true)
	assert.Equal(s.T(), trip.(models.GroceryTrip).CopyRemainingItems, false)
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

	s.mock.ExpectQuery("^SELECT count*").
		WithArgs(storeID).
		WillReturnRows(s.mock.NewRows([]string{"count"}).AddRow(1))

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

	trip, err := UpdateTrip(args)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), trip.(models.GroceryTrip).Completed, true)
	assert.Equal(s.T(), trip.(models.GroceryTrip).CopyRemainingItems, false)
}

func (s *Suite) TestUpdateTrip_MarkCompletedAndCopyRemainingItems() {
	tripID := uuid.NewV4()
	storeID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(tripID).
		WillReturnRows(s.mock.NewRows([]string{"id", "store_id", "name"}).AddRow(tripID, storeID, "Trip 1"))

	s.mock.ExpectExec("^UPDATE \"grocery_trips\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))

	s.mock.ExpectQuery("^SELECT count*").
		WithArgs(storeID).
		WillReturnRows(s.mock.NewRows([]string{"count"}).AddRow(1))

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
	s.mock.ExpectQuery("^INSERT INTO \"items\" (.+)$").
		WithArgs(newTripID, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 1, false, 1, AnyTime{}, AnyTime{}, nil).
		WillReturnRows(s.mock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))
	s.mock.ExpectExec("^UPDATE \"items\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))

	args := map[string]interface{}{
		"tripId":             tripID,
		"completed":          true,
		"copyRemainingItems": true,
	}
	trip, err := UpdateTrip(args)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), trip.(models.GroceryTrip).Completed, true)
	assert.Equal(s.T(), trip.(models.GroceryTrip).CopyRemainingItems, true)
}
