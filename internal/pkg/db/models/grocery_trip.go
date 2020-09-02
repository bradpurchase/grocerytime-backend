package models

import (
	"fmt"
	"time"

	// Postgres dialect for GORM
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	uuid "github.com/satori/go.uuid"
)

type GroceryTrip struct {
	ID                 uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	StoreID            uuid.UUID `gorm:"type:uuid;not null"`
	Name               string    `gorm:"type:varchar(100);not null"`
	Completed          bool      `gorm:"default:false;not null"`
	CopyRemainingItems bool      `gorm:"default:false;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	// Associations
	Store Store
}

// AfterUpdate hook is triggered after a trip is updated, such as in trips.UpdateTrip
func (g *GroceryTrip) AfterUpdate(tx *gorm.DB) (err error) {
	if g.Completed {
		var tripsCount int
		if err := tx.Model(&GroceryTrip{}).Where("store_id = ?", g.StoreID).Count(&tripsCount).Error; err != nil {
			return err
		}

		// Create the next trip for the user.
		// The new trip name is suffixed by a number that represents the number
		// of trips the user has made to the store (i.e. "Trip 12")
		newTripName := fmt.Sprintf("Trip %d", (tripsCount + 1))
		newTrip := GroceryTrip{StoreID: g.StoreID, Name: newTripName}
		if err := tx.Create(&newTrip).Error; err != nil {
			return err
		}

		// If the completed trip was configured to copy its remaining items
		// over to the next trip, perform this operation - otherwise, mark
		// each item in the completed trip as completed
		completed := true
		columns := Item{Completed: &completed}
		if g.CopyRemainingItems {
			columns = Item{GroceryTripID: newTrip.ID}
		}
		updateItemsQuery := tx.
			Model(&Item{}).
			Where("grocery_trip_id = ? AND completed = ?", g.ID, false).
			UpdateColumns(columns).
			Error
		if err := updateItemsQuery; err != nil {
			return err
		}
	}
	return nil
}
