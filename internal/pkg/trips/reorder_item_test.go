package trips

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestReorderItem_ReorderItemPosition(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: dbMock}), &gorm.Config{})
	require.NoError(t, err)

	itemID := uuid.NewV4()
	tripID := uuid.NewV4()

	mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(itemID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "grocery_trip_id"}).AddRow(itemID, tripID))

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
