package grocerylist

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
	listID := uuid.NewV4()
	userID := uuid.NewV4()
	itemRows := sqlmock.
		NewRows([]string{
			"id",
			"list_id",
			"user_id",
			"name",
			"quantity",
			"completed",
			"position",
			"created_at",
			"updated_at",
		}).
		AddRow(uuid.NewV4(), listID, userID, "Apples", 5, false, 1, time.Now(), time.Now()).
		AddRow(uuid.NewV4(), listID, userID, "Bananas", 10, false, 2, time.Now(), time.Now()).
		AddRow(itemID, listID, userID, "Oranges", 1, false, 3, time.Now(), time.Now())

	mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(itemID).
		WillReturnRows(itemRows)

	mock.ExpectBegin()
	mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(itemID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(itemID))
	mock.ExpectExec("^UPDATE \"items\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("^UPDATE \"lists\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectQuery("^SELECT (.+) FROM \"lists\"*").
		WithArgs(listID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(listID))

	list, err := ReorderItem(db, itemID, 4)
	require.NoError(t, err)
	assert.Equal(t, list.ID, listID)
}
