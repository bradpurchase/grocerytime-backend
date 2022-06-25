package meals

import (
	"github.com/DATA-DOG/go-sqlmock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *Suite) TestUpdateMeal_MealNotFound() {
	mealID := uuid.NewV4()
	args := map[string]interface{}{
		"id":       mealID,
		"name":     "PB&J",
		"mealType": "Lunch",
		"servings": 1,
		"date":     "2020-12-30",
	}

	s.mock.ExpectQuery("^SELECT (.+) FROM \"meals\"*").
		WithArgs(mealID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, e := UpdateMeal(args)
	require.Error(s.T(), e)
	assert.Equal(s.T(), "record not found", e.Error())
}

func (s *Suite) TestUpdateMeal_SingleColumn() {
	mealID := uuid.NewV4()
	name := "Peanut Butter & Jelly Sandwich"
	args := map[string]interface{}{
		"id":       mealID,
		"name":     name,
		"mealType": "Lunch",
		"servings": 1,
		"date":     "2020-12-30",
	}

	s.mock.ExpectQuery("^SELECT (.+) FROM \"meals\"*").
		WithArgs(mealID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(mealID, "PB&J"))

	meal, err := UpdateMeal(args)
	require.Error(s.T(), err)
	assert.Equal(s.T(), name, meal.Name)
}

func (s *Suite) TestUpdateMeal_MultiColumn() {
	mealID := uuid.NewV4()
	name := "Peanut Butter & Jelly Sandwich"
	mealType := "Snack"
	args := map[string]interface{}{
		"id":       mealID,
		"name":     name,
		"mealType": mealType,
		"servings": 1,
		"date":     "2020-12-30",
	}

	s.mock.ExpectQuery("^SELECT (.+) FROM \"meals\"*").
		WithArgs(mealID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "meal_type"}).AddRow(mealID, "PB&J", "Lunch"))

	meal, err := UpdateMeal(args)
	require.Error(s.T(), err)
	assert.Equal(s.T(), name, meal.Name)
	assert.Equal(s.T(), &mealType, meal.MealType)
}
