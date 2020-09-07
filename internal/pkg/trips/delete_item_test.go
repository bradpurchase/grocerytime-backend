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

func TestDeleteItem_ItemNotFound(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	itemID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(itemID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, e := DeleteItem(db, itemID)
	require.Error(t, e)
	assert.Equal(t, e.Error(), "item not found")
}

func TestDeleteItem_SuccessMoreInCategory(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	itemID := uuid.NewV4()
	categoryID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(itemID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "category_id"}).AddRow(itemID, categoryID))

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE \"items\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("^UPDATE \"grocery_trips\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectQuery("^SELECT count*").
		WithArgs(categoryID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	item, err := DeleteItem(db, itemID)
	require.NoError(t, err)
	assert.Equal(t, item.(models.Item).ID, itemID)
}

func TestDeleteItem_SuccessLastInCategory(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	itemID := uuid.NewV4()
	categoryID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(itemID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "category_id"}).AddRow(itemID, categoryID))

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE \"items\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("^UPDATE \"grocery_trips\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectQuery("^SELECT count*").
		WithArgs(categoryID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE \"grocery_trip_categories\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	item, err := DeleteItem(db, itemID)
	require.NoError(t, err)
	assert.Equal(t, item.(models.Item).ID, itemID)
}
