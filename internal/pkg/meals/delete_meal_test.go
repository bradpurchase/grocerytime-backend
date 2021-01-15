package meals

import (
	"github.com/DATA-DOG/go-sqlmock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *Suite) TestDeleteMeal_MealNotFound() {
	mealID := uuid.NewV4()
	userID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"meals\"*").
		WithArgs(mealID, userID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, e := DeleteMeal(mealID, userID)
	require.Error(s.T(), e)
	assert.Equal(s.T(), e.Error(), "record not found")
}

func (s *Suite) TestDeleteMeal_MealFound() {
	mealID := uuid.NewV4()
	userID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"meals\"*").
		WithArgs(mealID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(mealID, userID))

	s.mock.ExpectExec("^UPDATE \"meals\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectExec("^UPDATE \"meal_users\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectExec("^UPDATE \"items\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(0, 0))

	meal, err := DeleteMeal(mealID, userID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), meal.ID, mealID)
	assert.NotNil(s.T(), meal.DeletedAt)
}
