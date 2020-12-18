package meals

import (
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *Suite) TestCreateRecipe_FullDetails() {
	userID := uuid.NewV4()
	name := "PB&J"
	mealType := "Snack"
	url := "https://www.food.com/recipe/traditional-peanut-butter-and-jelly-243965"
	amount := 2.0
	var unit string
	unit1 := "tbsp"
	unit2 := "tsp"
	args := map[string]interface{}{
		"name":     name,
		"mealType": mealType,
		"url":      url,
		"ingredients": []interface{}{
			map[string]interface{}{
				"name":   "Bread",
				"amount": amount,
				"notes":  "with crusts cut off!",
			},
			map[string]interface{}{
				"name":   "Peanut Butter",
				"amount": amount,
				"unit":   unit1,
				"notes":  "spread evenly",
			},
			map[string]interface{}{
				"name":   "Strawberry Jam",
				"amount": amount,
				"unit":   unit2,
				"notes":  "spread evenly",
			},
		},
	}
	recipe, e := CreateRecipe(userID, args)
	require.Error(s.T(), e)
	assert.Equal(s.T(), recipe.Name, name)
	assert.Equal(s.T(), recipe.MealType, &mealType)
	assert.Equal(s.T(), recipe.URL, &url)

	assert.Equal(s.T(), len(recipe.Ingredients), 3)
	assert.Equal(s.T(), recipe.Ingredients[0].Name, "Bread")
	assert.Equal(s.T(), recipe.Ingredients[0].Amount, &amount)
	assert.Equal(s.T(), recipe.Ingredients[0].Unit, &unit)
	assert.Equal(s.T(), recipe.Ingredients[1].Name, "Peanut Butter")
	assert.Equal(s.T(), recipe.Ingredients[1].Amount, &amount)
	assert.Equal(s.T(), recipe.Ingredients[1].Unit, &unit1)
	assert.Equal(s.T(), recipe.Ingredients[2].Name, "Strawberry Jam")
	assert.Equal(s.T(), recipe.Ingredients[2].Amount, &amount)
	assert.Equal(s.T(), recipe.Ingredients[2].Unit, &unit2)
}
