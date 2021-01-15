package trips

import (
	"github.com/DATA-DOG/go-sqlmock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *Suite) TestDeleteItem_ItemNotFound() {
	itemID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(itemID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, e := DeleteItem(itemID)
	require.Error(s.T(), e)
	assert.Equal(s.T(), e.Error(), "item not found")
}

func (s *Suite) TestDeleteItem_SuccessMoreInCategory() {
	itemID := uuid.NewV4()
	categoryID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(itemID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "category_id"}).AddRow(itemID, categoryID))

	s.mock.ExpectBegin()
	s.mock.ExpectExec("^UPDATE \"items\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	s.mock.ExpectExec("^UPDATE \"grocery_trips\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))

	s.mock.ExpectQuery("^SELECT count*").
		WithArgs(categoryID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	item, err := DeleteItem(itemID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), item.ID, itemID)
}

func (s *Suite) TestDeleteItem_SuccessLastInCategory() {
	itemID := uuid.NewV4()
	categoryID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(itemID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "category_id"}).AddRow(itemID, categoryID))

	s.mock.ExpectBegin()
	s.mock.ExpectExec("^UPDATE \"items\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	s.mock.ExpectExec("^UPDATE \"grocery_trips\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))

	s.mock.ExpectQuery("^SELECT count*").
		WithArgs(categoryID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	s.mock.ExpectBegin()
	s.mock.ExpectExec("^UPDATE \"grocery_trip_categories\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	item, err := DeleteItem(itemID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), item.ID, itemID)
}
