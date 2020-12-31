package meals

import (
	"github.com/DATA-DOG/go-sqlmock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *Suite) TestPlanMeal_InvalidRecipeID() {
	userID := uuid.NewV4()
	storeID := uuid.NewV4()
	args := map[string]interface{}{
		"recipeId": "invalid",
		"storeId":  storeID.String(),
		"name":     "PB&J",
		"mealType": "Lunch",
		"servings": 1,
		"date":     "2020-12-30",
		"items": []interface{}{
			map[string]interface{}{
				"name":     "Peanut Butter",
				"quantity": 1,
			},
		},
	}
	_, e := PlanMeal(userID, args)
	require.Error(s.T(), e)
	assert.Equal(s.T(), e.Error(), "recipeId arg not a UUID")
}

func (s *Suite) TestPlanMeal_InvalidStoreID() {
	userID := uuid.NewV4()
	recipeID := uuid.NewV4()
	args := map[string]interface{}{
		"recipeId": recipeID.String(),
		"storeId":  "invalid",
		"name":     "PB&J",
		"mealType": "Lunch",
		"servings": 1,
		"date":     "2020-12-30",
		"items": []interface{}{
			map[string]interface{}{
				"name":     "Peanut Butter",
				"quantity": 1,
			},
		},
	}
	_, e := PlanMeal(userID, args)
	require.Error(s.T(), e)
	assert.Equal(s.T(), e.Error(), "storeId arg not a UUID")
}

func (s *Suite) TestPlanMeal_Valid() {
	userID := uuid.NewV4()
	storeID := uuid.NewV4()
	recipeID := uuid.NewV4()
	args := map[string]interface{}{
		"recipeId": recipeID.String(),
		"storeId":  storeID.String(),
		"name":     "PB&J",
		"mealType": "Snack",
		"servings": 1,
		"date":     "2020-12-30",
		"notes":    "with the crusts cut off!",
		"items": []interface{}{
			map[string]interface{}{
				"name":     "Peanut Butter",
				"quantity": 1,
			},
		},
	}

	mealID := uuid.NewV4()
	s.mock.ExpectBegin()

	s.mock.ExpectQuery("^INSERT INTO \"meals\" (.+)$").
		WithArgs(recipeID, userID, args["name"], args["mealType"], args["servings"], args["notes"], args["date"], AnyTime{}, AnyTime{}, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(mealID))

	storeName := "Test Store"
	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(storeID, storeName))
	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(userID, storeName).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name"}).AddRow(storeID, userID, storeName))
	tripID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(storeID, false).
		WillReturnRows(sqlmock.NewRows([]string{"id", "store_id"}).AddRow(tripID, storeID))
	s.mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(tripID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "store_id"}).AddRow(tripID, storeID))
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "store_id", "user_id"}).AddRow(uuid.NewV4(), storeID, userID))

	s.mock.ExpectCommit()

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