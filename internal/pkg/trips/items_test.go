package trips

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRetrieveItems_NoItems(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	tripID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(tripID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	items, err := RetrieveItems(db, tripID)
	require.NoError(t, err)
	assert.Equal(t, len(items.([]models.Item)), 0)
}

func TestRetrieveItems_HasItems(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	tripID := uuid.NewV4()
	itemRows := sqlmock.
		NewRows([]string{
			"id",
			"grocery_trip_id",
			"name",
			"quantity",
			"completed",
			"created_at",
			"updated_at",
		}).
		AddRow(uuid.NewV4(), tripID, "Apples", 5, false, time.Now(), time.Now()).
		AddRow(uuid.NewV4(), tripID, "Bananas", 2, false, time.Now(), time.Now())
	mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(tripID).
		WillReturnRows(itemRows)

	items, err := RetrieveItems(db, tripID)
	require.NoError(t, err)
	assert.Equal(t, len(items.([]models.Item)), 2)
	assert.Equal(t, items.([]models.Item)[0].GroceryTripID, tripID)
	assert.Equal(t, items.([]models.Item)[0].Name, "Apples")
	assert.Equal(t, items.([]models.Item)[1].Name, "Bananas")
}

func TestRetrieveItemsInCategory_NoItems(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	tripID := uuid.NewV4()
	categoryID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(tripID, categoryID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	items, err := RetrieveItemsInCategory(db, tripID, categoryID)
	require.NoError(t, err)
	assert.Equal(t, len(items.([]models.Item)), 0)
}

func TestRetrieveItemsInCategory_HasItems(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	tripID := uuid.NewV4()
	categoryID := uuid.NewV4()
	itemRows := sqlmock.
		NewRows([]string{
			"id",
			"grocery_trip_id",
			"name",
			"quantity",
			"completed",
			"created_at",
			"updated_at",
		}).
		AddRow(uuid.NewV4(), tripID, "Apples", 5, false, time.Now(), time.Now()).
		AddRow(uuid.NewV4(), tripID, "Bananas", 2, false, time.Now(), time.Now())
	mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(tripID, categoryID).
		WillReturnRows(itemRows)

	items, err := RetrieveItemsInCategory(db, tripID, categoryID)
	require.NoError(t, err)
	assert.Equal(t, len(items.([]models.Item)), 2)
	assert.Equal(t, items.([]models.Item)[0].GroceryTripID, tripID)
	assert.Equal(t, items.([]models.Item)[0].Name, "Apples")
	assert.Equal(t, items.([]models.Item)[1].Name, "Bananas")
}
