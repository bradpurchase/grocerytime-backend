package meals

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/grsmv/goweek"
	uuid "github.com/satori/go.uuid"
)

// PlannedMeals fetches the planned meals within the provided week/year for the current user
// Note: when used via the plannedMeals query, it is possible to leave the weekNumber/year
// nil and it will use the current time
func PlannedMeals(userID uuid.UUID, weekNumber int, year int) (meals []models.Meal, err error) {
	week, err := goweek.NewWeek(year, weekNumber)
	if err != nil {
		return meals, err
	}
	days := week.Days
	weekFirstDay := days[0].Format("2006-01-02")
	weekLastDay := days[len(days)-1].Format("2006-01-02")

	query := db.Manager.
		Preload("Users").
		Select("meals.*").
		Joins("INNER JOIN meal_users ON meal_users.meal_id = meals.id").
		Where("meal_users.user_id = ?", userID).
		Where("CAST(meals.created_at AS date) BETWEEN ? AND ?", weekFirstDay, weekLastDay).
		Order("CASE WHEN meals.meal_type='Breakfast' THEN 1 WHEN meals.meal_type='Lunch' THEN 2 WHEN meals.meal_type='Dinner' THEN 3 WHEN meals.meal_type='Dinner' THEN 4 WHEN meals.meal_type='Dessert' THEN 5 WHEN meals.meal_type='Snack' THEN 6 END ASC,meals.created_at DESC").
		Find(&meals).
		Error
	if err := query; err != nil {
		return meals, err
	}

	return meals, nil
}
