package models

import (
	"time"

	// Postgres dialect for GORM
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	uuid "github.com/satori/go.uuid"
)

type GroceryTrip struct {
	ID        uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	ListID    uuid.UUID `gorm:"type:uuid;not null"`
	Name      string    `gorm:"type:varchar(100);not null"`
	Completed bool      `gorm:"default:false;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time

	// Associations
	List List
}

// AfterUpdate hook is triggered after a trip is updated, such as in trips.UpdateTrip
func (g *GroceryTrip) AfterUpdate(tx *gorm.DB) (err error) {
	if g.Completed {
		// Update each item in a completed trip
		updateItemsQuery := tx.
			Model(&Item{}).
			Where("trip_id = ? AND completed = ?", g.ID, false).
			UpdateColumn("completed", true).
			Error
		if err := updateItemsQuery; err != nil {
			return err
		}

		// Create the next trip for the user
		newTrip := GroceryTrip{ListID: g.ListID, Name: "New Trip"}
		if err := tx.Create(&newTrip).Error; err != nil {
			return err
		}
	}
	return nil
}
