package stores

import (
	"github.com/DATA-DOG/go-sqlmock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *Suite) TestSaveStapleItem_ItemNotFound() {
	itemID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(itemID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	storeID := uuid.NewV4()
	_, err := SaveStapleItem(storeID, itemID)
	require.Error(s.T(), err)
	assert.Equal(s.T(), err.Error(), "record not found")
}

func (s *Suite) TestSaveStapleItem_FindExistingStapleItem() {
	itemID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(itemID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(itemID))

	stapleItemID := uuid.NewV4()
	storeID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_staple_items\"*").
		WithArgs(storeID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "store_id"}).AddRow(stapleItemID, storeID))

	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(itemID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(itemID))
	s.mock.ExpectExec("^UPDATE \"items\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))

	stapleItem, err := SaveStapleItem(storeID, itemID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), stapleItem.ID, stapleItemID)
}

func (s *Suite) TestSaveStapleItem_CreateNewStapleItem() {
	itemID := uuid.NewV4()
	itemName := "Apples"
	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(itemID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(itemID, itemName))

	stapleItemID := uuid.NewV4()
	storeID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_staple_items\"*").
		WithArgs(storeID, itemName).
		WillReturnRows(sqlmock.NewRows([]string{}))
	s.mock.ExpectQuery("^INSERT INTO \"store_staple_items\" (.+)$").
		WithArgs(storeID, itemName, AnyTime{}, AnyTime{}, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(stapleItemID))

	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(itemID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(itemID))
	s.mock.ExpectExec("^UPDATE \"items\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))

	stapleItem, err := SaveStapleItem(storeID, itemID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), stapleItem.ID, stapleItemID)
}
