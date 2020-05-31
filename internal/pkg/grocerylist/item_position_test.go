package grocerylist

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetermineListPosition_TopOfList(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	listID := uuid.NewV4()

	itemRows := sqlmock.
		NewRows([]string{
			"list_id",
			"position",
		}).
		AddRow(listID, 1000).
		AddRow(listID, 998)
	mock.ExpectQuery("^SELECT position FROM \"items\"*").
		WithArgs(listID).
		WillReturnRows(itemRows)

	pos, err := DetermineListPosition("top", db, listID)
	require.NoError(t, err)
	assert.Equal(t, pos, 996)
}

func TestDetermineListPosition_BottomOfList(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	listID := uuid.NewV4()

	itemRows := sqlmock.
		NewRows([]string{
			"list_id",
			"position",
		}).
		AddRow(listID, 998).
		AddRow(listID, 1000)
	mock.ExpectQuery("^SELECT position FROM \"items\"*").
		WithArgs(listID).
		WillReturnRows(itemRows)

	pos, err := DetermineListPosition("bottom", db, listID)
	require.NoError(t, err)
	assert.Equal(t, pos, 1002)
}
