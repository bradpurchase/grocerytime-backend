package meals

import (
	"github.com/DATA-DOG/go-sqlmock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *Suite) TestPlanMeal_InvalidRecipeID() {
	userID := uuid.NewV4()
	args := map[string]interface{}{
		"recipeId": "invalid",
		"name":     "PB&J",
		"mealType": "Lunch",
		"servings": 1,
		"date":     "2020-12-30",
	}
	_, e := PlanMeal(userID, args)
	require.Error(s.T(), e)
	assert.Equal(s.T(), e.Error(), "recipeId arg not a UUID")
}

func (s *Suite) TestPlanMeal_Valid() {
	userID := uuid.NewV4()
	recipeID := uuid.NewV4()
	args := map[string]interface{}{
		"recipeId": recipeID.String(),
		"name":     "PB&J",
		"mealType": "Snack",
		"servings": 1,
		"date":     "2020-12-30",
		"notes":    "with the crusts cut off!",
	}

	mealID := uuid.NewV4()
	s.mock.ExpectQuery("^INSERT INTO \"meals\" (.+)$").
		WithArgs(recipeID, userID, args["name"], args["mealType"], args["servings"], args["notes"], args["date"], AnyTime{}, AnyTime{}, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(mealID))

	meal, err := PlanMeal(userID, args)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), meal.ID, mealID)
	assert.Equal(s.T(), meal.RecipeID, recipeID)
	assert.Equal(s.T(), meal.UserID, userID)
	assert.Equal(s.T(), meal.Name, args["name"])
	mealType := args["mealType"].(string)
	assert.Equal(s.T(), meal.MealType, &mealType)
	assert.Equal(s.T(), meal.Servings, args["servings"])
	assert.Equal(s.T(), meal.Date, args["date"])
}
