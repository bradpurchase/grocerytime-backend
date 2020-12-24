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

	recipes, err := RetrieveRecipes(userID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), len(recipes), 0)
}

func (s *Suite) TestRetrieveRecipes_WithRecipes() {
	userID := uuid.NewV4()
	recipeID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"recipes\"*").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(recipeID, userID))

	recipes, err := RetrieveRecipes(userID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), len(recipes), 1)
	assert.Equal(s.T(), recipes[0].UserID, userID)
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
