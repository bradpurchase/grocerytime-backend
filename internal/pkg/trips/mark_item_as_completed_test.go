package trips

import (
	"github.com/DATA-DOG/go-sqlmock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *Suite) TestMarkItemAsCompleted_CouldNotUpdate() {
	userID := uuid.NewV4()
	_, err := MarkItemAsCompleted("", userID)
	require.Error(s.T(), err)
	assert.Equal(s.T(), err.Error(), "could not update items")
}

func (s *Suite) TestMarkItemAsCompleted_Updated() {
	userID := uuid.NewV4()
	name := "Apples"

	s.mock.ExpectExec("^UPDATE \"items\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	itemID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(name, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(itemID))

	_, err := MarkItemAsCompleted(name, userID)
	require.NoError(s.T(), err)
}
