package stores

import (
	"github.com/DATA-DOG/go-sqlmock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *Suite) TestRemoveStapleItem_StapleItemNotFound() {
	itemID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(itemID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, err := RemoveStapleItem(itemID)
	require.Error(s.T(), err)
	assert.Equal(s.T(), err.Error(), "record not found")
}

func (s *Suite) TestRemoveStapleItem_StapleItemRemoved() {
	itemID := uuid.NewV4()
	stapleItemID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(itemID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "staple_item_id"}).AddRow(itemID, stapleItemID))

	s.mock.ExpectExec("^UPDATE \"store_staple_items\"*").
		WillReturnResult(sqlmock.NewResult(1, 1))

	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(itemID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(itemID))
	s.mock.ExpectExec("^UPDATE \"items\"*").
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err := RemoveStapleItem(itemID)
	require.NoError(s.T(), err)
}
