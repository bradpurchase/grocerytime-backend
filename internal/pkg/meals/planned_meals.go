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
	weekFirstDay := days[0]
	weekLastDay := days[len(days)-1]

	// TODO: inner join with meal_users and select user_id that way
	query := db.Manager.
		Where("user_id = ?", userID).
		Where("created_at BETWEEN ? AND ?", weekFirstDay, weekLastDay).
		Find(&meals).
		Order("created_at DESC").
		Error
	if err := query; err != nil {
		return meals, err
	}

	return meals, nil
}
