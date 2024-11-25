package meals

import (
	"time"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/grsmv/goweek"
	uuid "github.com/satori/go.uuid"
)

// RetrieveMeals fetches the planned meals within the provided week/year for the current user
// Note: when used via the plannedMeals query, it is possible to leave the weekNumber/year
// nil and it will use the current time
func RetrieveMeals(userID uuid.UUID, args map[string]interface{}) (meals []models.Meal, err error) {
	query := db.Manager.
		Preload("Users").
		Select("meals.*").
		Joins("INNER JOIN meal_users ON meal_users.meal_id = meals.id").
		Where("meal_users.user_id = ?", userID)

	// Support fetching meals by date range
	// we allow no value for year and weekNumber args and use current time for this case
	year := args["year"]
	weekNumber := args["weekNumber"]
	if year == nil || weekNumber == nil {
		t := time.Now()
		year, weekNumber = t.ISOWeek()
	}
	week, err := goweek.NewWeek(year.(int), weekNumber.(int))
	if err != nil {
		return meals, err
	}
	days := week.Days
	weekFirstDay := days[0].Format("2006-01-02")
	weekLastDay := days[len(days)-1].Format("2006-01-02")
	query = query.Where("CAST(meals.date AS date) BETWEEN ? AND ?", weekFirstDay, weekLastDay)

	// Allow filtering by meal_type
	if args["mealType"] != nil && args["mealType"] != "" {
		query = query.Where("meal_type = ?", args["mealType"].(string))
	}

	query = query.
		Order("CASE WHEN meals.meal_type='Breakfast' THEN 1 WHEN meals.meal_type='Lunch' THEN 2 WHEN meals.meal_type='Dinner' THEN 3 WHEN meals.meal_type='Dinner' THEN 4 WHEN meals.meal_type='Dessert' THEN 5 WHEN meals.meal_type='Snack' THEN 6 END ASC").
		Order("meals.date DESC").
		Find(&meals)
	if err := query.Error; err != nil {
		return meals, err
	}

	return meals, nil
}

// RetrieveMealForUser retrieves a specific meal by mealID and userID
func RetrieveMealForUser(mealID uuid.UUID, userID uuid.UUID) (meal models.Meal, err error) {
	query := db.Manager.
		Select("meals.*").
		Joins("INNER JOIN meal_users ON meal_users.meal_id = meals.id").
		Where("meals.id = ?", mealID).
		Where("meal_users.user_id = ?", userID).
		Last(&meal).
		Error
	if err := query; err != nil {
		return meal, err
	}
	return meal, nil
}
