package trips

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReorder_ReorderItemPosition(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	itemID := uuid.NewV4()
	tripID := uuid.NewV4()
	userID := uuid.NewV4()
	itemRows := sqlmock.
		NewRows([]string{
			"id",
			"grocery_trip_id",
			"user_id",
			"name",
			"quantity",
			"completed",
			"position",
			"created_at",
			"updated_at",
		}).
		AddRow(uuid.NewV4(), tripID, userID, "Apples", 5, false, 1, time.Now(), time.Now()).
		AddRow(uuid.NewV4(), tripID, userID, "Bananas", 10, false, 2, time.Now(), time.Now()).
		AddRow(itemID, tripID, userID, "Oranges", 1, false, 3, time.Now(), time.Now())

	mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(itemID).
		WillReturnRows(itemRows)

	mock.ExpectBegin()
	mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(itemID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(itemID))
	mock.ExpectExec("^UPDATE \"items\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("^UPDATE \"grocery_trips\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(tripID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tripID))

	trip, err := ReorderItem(db, itemID, 4)
	require.NoError(t, err)
	assert.Equal(t, trip.ID, tripID)
}
