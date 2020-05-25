package grocerylist

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

func TestRetrieveItemsInList_NoItems(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	listID := uuid.NewV4()
	list := &models.List{ID: listID, Name: "Test List", CreatedAt: time.Now(), UpdatedAt: time.Now()}

	mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(listID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	items, err := RetrieveItemsInList(db, list.ID)
	require.NoError(t, err)
	assert.Equal(t, len(items.([]models.Item)), 0)
}

func TestRetrieveItemsInList_HasItems(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	listID := uuid.NewV4()
	list := &models.List{ID: listID, Name: "Test List", CreatedAt: time.Now(), UpdatedAt: time.Now()}

	itemRows := sqlmock.
		NewRows([]string{
			"id",
			"list_id",
			"name",
			"quantity",
			"completed",
			"created_at",
			"updated_at",
		}).
		AddRow(uuid.NewV4(), listID, "Apples", 5, false, time.Now(), time.Now()).
		AddRow(uuid.NewV4(), listID, "Bananas", 2, false, time.Now(), time.Now())
	mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(listID).
		WillReturnRows(itemRows)

	items, err := RetrieveItemsInList(db, list.ID)
	require.NoError(t, err)
	assert.Equal(t, len(items.([]models.Item)), 2)
	assert.Equal(t, items.([]models.Item)[0].ListID, list.ID)
	assert.Equal(t, items.([]models.Item)[0].Name, "Apples")
	assert.Equal(t, items.([]models.Item)[1].Name, "Bananas")
}

func TestRetrieveItemsInList_Order(t *testing.T) {
	//TODO write test proving order.. new ones on top, completed ones on bottom
}
