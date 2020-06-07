package trips

import (
	"database/sql/driver"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestUpdateTrip_TripNotFound(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	tripID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(tripID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	args := map[string]interface{}{
		"tripId": tripID,
	}

	_, e := UpdateTrip(db, args)
	require.Error(t, e)
	assert.Equal(t, e.Error(), "record not found")
}

func TestUpdateTrip_NameUpdate(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	tripID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(tripID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(tripID, "My First Trip"))

	args := map[string]interface{}{
		"tripId": tripID,
		"name":   "My Second Trip",
	}

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE \"grocery_trips\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	trip, err := UpdateTrip(db, args)
	require.NoError(t, err)
	assert.Equal(t, trip.(models.GroceryTrip).Name, "My Second Trip")
}

func TestUpdateTrip_MarkCompleted(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	tripID := uuid.NewV4()
	listID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(tripID).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "list_id", "name"}).
			AddRow(tripID, listID, "My First Trip"))

	args := map[string]interface{}{
		"tripId":    tripID,
		"completed": true,
	}

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE \"grocery_trips\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery("^INSERT INTO \"grocery_trips\" (.+)$").
		WithArgs(listID, "New Trip", AnyTime{}, AnyTime{}).
		WillReturnRows(sqlmock.NewRows([]string{"list_id"}).AddRow(listID))
	mock.ExpectExec("^UPDATE \"items\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	trip, err := UpdateTrip(db, args)
	require.NoError(t, err)
	assert.Equal(t, trip.(models.GroceryTrip).Completed, true)
}
