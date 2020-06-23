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

func TestRetrieveCurrentTripInList_UserNotAMemberOfList(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	listID := uuid.NewV4()
	user := models.User{ID: uuid.NewV4()}
	mock.ExpectQuery("^SELECT (.+) FROM \"list_users\"*").
		WithArgs(listID, user.ID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, e := RetrieveCurrentTripInList(db, listID, user)
	require.Error(t, e)
	assert.Equal(t, e.Error(), "user is not a member of this list")
}

func TestRetrieveCurrentTripInList_TripNotAssociatedWithList(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	listID := uuid.NewV4()
	user := models.User{ID: uuid.NewV4(), Email: "test@example.com"}
	mock.ExpectQuery("^SELECT (.+) FROM \"list_users\"*").
		WithArgs(listID, user.ID, user.Email).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))

	_, e := RetrieveCurrentTripInList(db, listID, user)
	require.Error(t, e)
	assert.Equal(t, e.Error(), "could not find trip associated with this list")
}

func TestRetrieveCurrentTripInList_FoundResult(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	listID := uuid.NewV4()
	user := models.User{ID: uuid.NewV4(), Email: "test@example.com"}
	mock.ExpectQuery("^SELECT (.+) FROM \"list_users\"*").
		WithArgs(listID, user.ID, user.Email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(uuid.NewV4(), user.ID))

	tripID := uuid.NewV4()
	tripName := "Week of May 31"
	mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(listID, false).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(tripID, tripName))

	trip, err := RetrieveCurrentTripInList(db, listID, user)
	require.NoError(t, err)
	assert.Equal(t, trip.(models.GroceryTrip).Name, tripName)
}
