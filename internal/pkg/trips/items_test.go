package trips

import (
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *Suite) TestRetrieveItems_NoItems() {
	tripID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(tripID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	items, err := RetrieveItems(tripID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), len(items.([]models.Item)), 0)
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
			"created_at",
			"updated_at",
		}).
		AddRow(uuid.NewV4(), tripID, "Apples", 5, false, time.Now(), time.Now()).
		AddRow(uuid.NewV4(), tripID, "Bananas", 2, false, time.Now(), time.Now())
	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(tripID).
		WillReturnRows(itemRows)

	items, err := RetrieveItems(tripID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), len(items.([]models.Item)), 2)
	assert.Equal(s.T(), items.([]models.Item)[0].GroceryTripID, tripID)
	assert.Equal(s.T(), items.([]models.Item)[0].Name, "Apples")
	assert.Equal(s.T(), items.([]models.Item)[1].Name, "Bananas")
}

func (s *Suite) TestRetrieveItemsInCategory_NoItems() {
	tripID := uuid.NewV4()
	categoryID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(tripID, categoryID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	items, err := RetrieveItemsInCategory(tripID, categoryID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), len(items.([]models.Item)), 0)
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
			"created_at",
			"updated_at",
		}).
		AddRow(uuid.NewV4(), tripID, "Apples", 5, false, time.Now(), time.Now()).
		AddRow(uuid.NewV4(), tripID, "Bananas", 2, false, time.Now(), time.Now())
	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(tripID, categoryID).
		WillReturnRows(itemRows)

	items, err := RetrieveItemsInCategory(tripID, categoryID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), len(items.([]models.Item)), 2)
	assert.Equal(s.T(), items.([]models.Item)[0].GroceryTripID, tripID)
	assert.Equal(s.T(), items.([]models.Item)[0].Name, "Apples")
	assert.Equal(s.T(), items.([]models.Item)[1].Name, "Bananas")
}
