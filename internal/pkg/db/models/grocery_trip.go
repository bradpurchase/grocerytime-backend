package models

import (
	"time"

	// Postgres dialect for GORM
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	uuid "github.com/satori/go.uuid"
)

type GroceryTrip struct {
	ID                 uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	ListID             uuid.UUID `gorm:"type:uuid;not null"`
	Name               string    `gorm:"type:varchar(100);not null"`
	Completed          bool      `gorm:"default:false;not null"`
	CopyRemainingItems bool      `gorm:"default:false;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time

	// Associations
	List List
}

// AfterUpdate hook is triggered after a trip is updated, such as in trips.UpdateTrip
func (g *GroceryTrip) AfterUpdate(tx *gorm.DB) (err error) {
	// Update lists.updated_at when a trip is updated
	// Note: trips are also updated when items are added/updated/created
	tx.Model(&List{}).Where("id = ?", g.ListID).Update("updated_at", time.Now())

	if g.Completed {
		// Create the next trip for the user
		newTrip := GroceryTrip{ListID: g.ListID, Name: "New Trip"}
		if err := tx.Create(&newTrip).Error; err != nil {
			return err
		}

		// If the completed trip was configured to copy its remaining items
		// over to the next trip, perform this operation - otherwise, mark
		// each item in the completed trip as completed
		columns := Item{Completed: true}
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
