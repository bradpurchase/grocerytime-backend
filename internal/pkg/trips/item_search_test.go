package trips

import (
	"fmt"

	"github.com/DATA-DOG/go-sqlmock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *Suite) TestSearchForItemByName_NoTripHistory() {
	userID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	_, e := SearchForItemByName("zap", userID)
	require.Error(s.T(), e)
	assert.Equal(s.T(), e.Error(), "no item matches the search term")
}

func (s *Suite) TestSearchForItemByName_ItemNotFound() {
	userID := uuid.NewV4()
	tripID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tripID))

	name := "zonk"
	nameArg := fmt.Sprintf("%%%s%%", name)
	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(nameArg, tripID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

	result, err := SearchForItemByName(name, userID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), result.Name, "")
}

func (s *Suite) TestSearchForItemByName_Found() {
	userID := uuid.NewV4()
	tripID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tripID))

	name := "app"
	nameArg := fmt.Sprintf("%%%s%%", name)
	rows := sqlmock.
		NewRows([]string{"id", "name"}).
		AddRow(uuid.NewV4(), "Apples")
	s.mock.ExpectQuery("^SELECT (.+) FROM \"items\"*").
		WithArgs(nameArg, tripID).
		WillReturnRows(rows)

	item, err := SearchForItemByName(name, userID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), item.Name, "Apples")
}
