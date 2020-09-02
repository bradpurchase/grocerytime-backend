package trips

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRetrieveCurrentStoreTrip_UserNotAMemberOfStore(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	storeID := uuid.NewV4()
	user := models.User{ID: uuid.NewV4()}
	mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID, user.ID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, e := RetrieveCurrentStoreTrip(db, storeID, user)
	require.Error(t, e)
	assert.Equal(t, e.Error(), "user is not a member of this store")
}

func TestRetrieveCurrentStoreTrip_TripNotAssociatedWithStore(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	storeID := uuid.NewV4()
	user := models.User{ID: uuid.NewV4(), Email: "test@example.com"}
	mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID, user.ID, user.Email).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))

	_, e := RetrieveCurrentStoreTrip(db, storeID, user)
	require.Error(t, e)
	assert.Equal(t, e.Error(), "could not find trip associated with this store")
}

func TestRetrieveCurrentStoreTrip_FoundResult(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	storeID := uuid.NewV4()
	user := models.User{ID: uuid.NewV4(), Email: "test@example.com"}
	mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID, user.ID, user.Email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(uuid.NewV4(), user.ID))

	tripID := uuid.NewV4()
	tripName := "Week of May 31"
	mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(storeID, false).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(tripID, tripName))

	trip, err := RetrieveCurrentStoreTrip(db, storeID, user)
	require.NoError(t, err)
	assert.Equal(t, trip.(models.GroceryTrip).Name, tripName)
}

func TestRetrieveTrip_NotFound(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	tripID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(tripID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, e := RetrieveTrip(db, tripID)
	require.Error(t, e)
	assert.Equal(t, e.Error(), "trip not found")
}

func TestRetrieveTrip_Found(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	tripID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(tripID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tripID))

	trip, err := RetrieveTrip(db, tripID)
	require.NoError(t, err)
	assert.Equal(t, trip.ID, tripID)
}
