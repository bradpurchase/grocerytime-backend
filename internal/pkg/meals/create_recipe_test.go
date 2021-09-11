package meals

import (
	"encoding/json"

	"github.com/DATA-DOG/go-sqlmock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"
)

func (s *Suite) TestCreateRecipe_NoIngredients() {
	userID := uuid.NewV4()
	var instructions []interface{}
	instructionsEncoded, _ := json.Marshal(instructions)
	instructionsJSON := datatypes.JSON(instructionsEncoded)

	args := map[string]interface{}{
		"name":        "PB&J",
		"mealType":    "Snack",
		"url":         "https://www.food.com/recipe/traditional-peanut-butter-and-jelly-243965",
		"ingredients": instructions,
	}

	recipeID := uuid.NewV4()
	s.mock.ExpectQuery("^INSERT INTO \"recipes\" (.+)$").
		WithArgs(sqlmock.AnyArg(), args["name"], "", instructionsJSON, args["mealType"], args["url"], "", AnyTime{}, AnyTime{}, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(recipeID))

	recipe, err := CreateRecipe(userID, args)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), recipe.ID, recipeID)
	assert.Equal(s.T(), recipe.Name, args["name"])
	mealType := args["mealType"].(string)
	assert.Equal(s.T(), recipe.MealType, &mealType)
	url := args["url"].(string)
	assert.Equal(s.T(), recipe.URL, &url)
	assert.Equal(s.T(), len(recipe.Ingredients), 0)
}

func (s *Suite) TestCreateRecipe_FullDetails() {
	userID := uuid.NewV4()

	// Ingredients
	ingName := "Bread"
	amount := "2.0"
	var unit string
	var notes string
	ingName1 := "Peanut Butter"
	unit1 := "tbsp"
	notes1 := "spread evenly"
	ingName2 := "Strawberry Jam"
	unit2 := "tsp"

	var instructions []interface{}

	args := map[string]interface{}{
		"name":        "PB&J",
		"description": "Nothing fancy, just a classic. Either smooth or crunch peanut butter is acceptable. Classically, the jelly is either strawberry or grape.",
		"mealType":    "Lunch",
		"url":         "https://www.food.com/recipe/traditional-peanut-butter-and-jelly-243965",
		"imageUrl":    "https://img.sndimg.com/food/image/upload/c_thumb,q_80,w_596,h_335/v1/img/recipes/24/39/65/picIDMFir.jpg",
		"ingredients": []interface{}{
			map[string]interface{}{
				"name":   ingName,
				"amount": amount,
			},
			map[string]interface{}{
				"name":   ingName1,
				"amount": amount,
				"unit":   unit1,
				"notes":  notes1,
			},
			map[string]interface{}{
				"name":   ingName2,
				"amount": amount,
				"unit":   unit2,
				"notes":  notes1,
			},
		},
		"instructions": instructions,
	}

	instructionsEncoded, _ := json.Marshal(instructions)
	instructionsJSON := datatypes.JSON(instructionsEncoded)

	recipeID := uuid.NewV4()
	s.mock.ExpectQuery("^INSERT INTO \"recipes\" (.+)$").
		WithArgs(sqlmock.AnyArg(), args["name"], args["description"], instructionsJSON, args["mealType"], args["url"], args["imageUrl"], AnyTime{}, AnyTime{}, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(recipeID))

	// Note: because we are creating the recipe_ingredients using the association,
	// GORM does a bulk insert - this is why there are many args here
	s.mock.ExpectQuery("^INSERT INTO \"recipe_ingredients\" (.+)$").
		WithArgs(recipeID, ingName, amount, unit, notes, AnyTime{}, AnyTime{}, nil, recipeID, ingName1, amount, unit1, notes1, AnyTime{}, AnyTime{}, nil, recipeID, ingName2, amount, unit2, notes1, AnyTime{}, AnyTime{}, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))

	recipe, err := CreateRecipe(userID, args)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), recipe.ID, recipeID)
	assert.Equal(s.T(), recipe.Name, args["name"])
	description := args["description"].(string)
	assert.Equal(s.T(), recipe.Description, &description)
	mealType := args["mealType"].(string)
	assert.Equal(s.T(), recipe.MealType, &mealType)
	url := args["url"].(string)
	assert.Equal(s.T(), recipe.URL, &url)
	imageURL := args["imageUrl"].(string)
	assert.Equal(s.T(), recipe.ImageURL, &imageURL)

	assert.Equal(s.T(), len(recipe.Ingredients), 3)
	assert.Equal(s.T(), recipe.Ingredients[0].Name, ingName)
	assert.Equal(s.T(), recipe.Ingredients[0].Amount, &amount)
	assert.Equal(s.T(), recipe.Ingredients[0].Unit, &unit)
	assert.Equal(s.T(), recipe.Ingredients[0].Notes, &notes)
	assert.Equal(s.T(), recipe.Ingredients[1].Name, ingName1)
	assert.Equal(s.T(), recipe.Ingredients[1].Amount, &amount)
	assert.Equal(s.T(), recipe.Ingredients[1].Unit, &unit1)
	assert.Equal(s.T(), recipe.Ingredients[1].Notes, &notes1)
	assert.Equal(s.T(), recipe.Ingredients[2].Name, ingName2)
	assert.Equal(s.T(), recipe.Ingredients[2].Amount, &amount)
	assert.Equal(s.T(), recipe.Ingredients[2].Unit, &unit2)
	assert.Equal(s.T(), recipe.Ingredients[2].Notes, &notes1)
}
