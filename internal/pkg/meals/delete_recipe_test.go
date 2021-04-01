package meals

import (
	"github.com/DATA-DOG/go-sqlmock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *Suite) TestDeleteRecipe_RecipeNotFound() {
	recipeID := uuid.NewV4()
	userID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"recipes\"*").
		WithArgs(recipeID, userID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, e := DeleteRecipe(recipeID, userID)
	require.Error(s.T(), e)
	assert.Equal(s.T(), e.Error(), "record not found")
}

func (s *Suite) TestDeleteRecipe_RecipeFound() {
	recipeID := uuid.NewV4()
	userID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"recipes\"*").
		WithArgs(recipeID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(recipeID))

	s.mock.ExpectExec("^UPDATE \"recipes\" SET (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))

	recipe, err := DeleteRecipe(recipeID, userID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), recipe.ID, recipeID)
	assert.NotNil(s.T(), recipe.DeletedAt)
}
