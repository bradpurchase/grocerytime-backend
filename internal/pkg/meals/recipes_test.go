package meals

import (
	"github.com/DATA-DOG/go-sqlmock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *Suite) TestRetrieveRecipes_NoRecipes() {
	userID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"recipes\"*").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	var args map[string]interface{}
	recipes, err := RetrieveRecipes(userID, args)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), len(recipes), 0)
}

func (s *Suite) TestRetrieveRecipes_RecipesNoMealTypeFilter() {
	userID := uuid.NewV4()
	recipeID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"recipes\"*").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(recipeID, userID))
	s.mock.ExpectQuery("^SELECT (.+) FROM \"recipe_ingredients\"*").
		WithArgs(recipeID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "recipe_id"}).AddRow(uuid.NewV4(), recipeID))

	var args map[string]interface{}
	recipes, err := RetrieveRecipes(userID, args)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), len(recipes), 1)
	assert.Equal(s.T(), recipes[0].UserID, userID)
}

func (s *Suite) TestRetrieveRecipes_RecipesWithMealTypeFilter() {
	userID := uuid.NewV4()
	recipeID := uuid.NewV4()
	mealTypeStr := "Breakfast"
	args := map[string]interface{}{
		"mealType": mealTypeStr,
	}
	rows := sqlmock.
		NewRows([]string{"id", "user_id", "name", "meal_type"}).
		AddRow(recipeID, userID, "Green Eggs & Ham", "Breakfast").
		AddRow(recipeID, userID, "Western Omelette", "Breakfast")
	s.mock.ExpectQuery("^SELECT (.+) FROM \"recipes\"*").
		WithArgs(userID, mealTypeStr).
		WillReturnRows(rows)
	s.mock.ExpectQuery("^SELECT (.+) FROM \"recipe_ingredients\"*").
		WithArgs(recipeID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "recipe_id"}).AddRow(uuid.NewV4(), recipeID))

	recipes, err := RetrieveRecipes(userID, args)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), len(recipes), 2)
	assert.Equal(s.T(), recipes[0].MealType, &mealTypeStr)
	assert.Equal(s.T(), recipes[1].MealType, &mealTypeStr)
}

func (s *Suite) TestRetrieveRecipes_RecipesWithLimit() {
	userID := uuid.NewV4()
	recipeID := uuid.NewV4()
	limit := 2
	args := map[string]interface{}{
		"limit": limit,
	}
	rows := sqlmock.
		NewRows([]string{"id", "user_id", "name", "meal_type"}).
		AddRow(recipeID, userID, "Green Eggs & Ham", "Breakfast").
		AddRow(recipeID, userID, "Western Omelette", "Breakfast")
	s.mock.ExpectQuery("^SELECT (.+) FROM \"recipes\" (.+) LIMIT 2*").
		WithArgs(userID).
		WillReturnRows(rows)
	s.mock.ExpectQuery("^SELECT (.+) FROM \"recipe_ingredients\"*").
		WithArgs(recipeID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "recipe_id"}).AddRow(uuid.NewV4(), recipeID))

	recipes, err := RetrieveRecipes(userID, args)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), len(recipes), limit)
}

func (s *Suite) TestRetrieveRecipe_NotFound() {
	recipeID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"recipes\"*").
		WithArgs(recipeID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, e := RetrieveRecipe(recipeID)
	require.Error(s.T(), e)
	assert.Equal(s.T(), e.Error(), "record not found")
}

func (s *Suite) TestRetrieveRecipe_Found() {
	recipeID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"recipes\"*").
		WithArgs(recipeID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(recipeID))
	s.mock.ExpectQuery("^SELECT (.+) FROM \"recipe_ingredients\"*").
		WithArgs(recipeID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "recipe_id"}).AddRow(uuid.NewV4(), recipeID))

	recipe, err := RetrieveRecipe(recipeID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), recipe.ID, recipeID)
}
