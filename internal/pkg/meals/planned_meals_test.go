package meals

import (
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *Suite) TestPlannedMeals_InvalidWeekNumber() {
	userID := uuid.NewV4()
	weekNumber := 54
	year := 2021
	_, e := PlannedMeals(userID, weekNumber, year)
	require.Error(s.T(), e)
	assert.Contains(s.T(), e.Error(), "number of week can't be less than 1 or greater than 53")
}

func (s *Suite) TestPlannedMeals_NoMeals() {
	userID := uuid.NewV4()
	weekNumber := 1
	year := 2021

	weekFirstDay := time.Date(2021, time.January, 04, 0, 0, 0, 0, time.UTC)
	weekLastDay := time.Date(2021, time.January, 10, 0, 0, 0, 0, time.UTC)
	s.mock.ExpectQuery("^SELECT meals.* FROM \"meals\"*").
		WithArgs(userID, weekFirstDay, weekLastDay).
		WillReturnRows(sqlmock.NewRows([]string{}))

	meals, err := PlannedMeals(userID, weekNumber, year)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), len(meals), 0)
}

func (s *Suite) TestPlannedMeals_ExistingMeals() {
	userID := uuid.NewV4()
	weekNumber := 1
	year := 2021

	weekFirstDay := time.Date(2021, time.January, 04, 0, 0, 0, 0, time.UTC)
	weekLastDay := time.Date(2021, time.January, 10, 0, 0, 0, 0, time.UTC)

	meal1ID := uuid.NewV4()
	meal2ID := uuid.NewV4()
	mealRows := sqlmock.NewRows([]string{"id"}).
		AddRow(meal1ID).
		AddRow(meal2ID)
	s.mock.ExpectQuery("^SELECT meals.* FROM \"meals\"*").
		WithArgs(userID, weekFirstDay, weekLastDay).
		WillReturnRows(mealRows)

	mealUserRows := sqlmock.NewRows([]string{"id", "meal_id", "user_id"}).
		AddRow(uuid.NewV4(), meal1ID, userID).
		AddRow(uuid.NewV4(), meal2ID, userID)
	s.mock.ExpectQuery("^SELECT (.+) FROM \"meal_users\"*").
		WithArgs(meal1ID, meal2ID).
		WillReturnRows(mealUserRows)

	meals, err := PlannedMeals(userID, weekNumber, year)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), len(meals), 2)
	assert.Equal(s.T(), meals[0].ID, meal1ID)
	assert.Equal(s.T(), meals[1].ID, meal2ID)
}
