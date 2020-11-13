package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// Item defines the model for items
type Item struct {
	ID            uuid.UUID  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	GroceryTripID uuid.UUID  `gorm:"type:uuid;not null;index:idx_items_grocery_trip_id_name"`
	CategoryID    *uuid.UUID `gorm:"type:uuid;not null"`
	UserID        uuid.UUID  `gorm:"type:uuid;not null"`
	Name          string     `gorm:"type:varchar(100);not null;index:idx_items_grocery_trip_id_name"`
	Quantity      int        `gorm:"default:1;not null"`
	Completed     *bool      `gorm:"default:false;not null"`
	Position      int        `gorm:"default:1;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	// Associations
	GroceryTrip GroceryTrip
}

func (i *Item) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Exec("UPDATE items SET position = position + 1 WHERE grocery_trip_id = ? AND position >= 0", i.GroceryTripID)
	return nil
}

// AfterCreate hook to touch the associated grocery trip after an item is created
// so that its UpdatedAt column is updated
func (i *Item) AfterCreate(tx *gorm.DB) (err error) {
	tx.Model(&GroceryTrip{}).Where("id = ?", i.GroceryTripID).Update("updated_at", time.Now())
	return nil
}

// BeforeUpdate hook handles reordering items
func (i *Item) BeforeUpdate(tx *gorm.DB) (err error) {
	item := &Item{}
	if err := tx.Where("id = ?", i.ID).Find(&item).Error; err != nil {
		return err
	}
	currPosition := item.Position
	newPosition := i.Position
	if currPosition == newPosition {
		return nil
	}
	if currPosition > newPosition {
		tx.Exec("UPDATE items SET position = position + 1 WHERE grocery_trip_id = ? AND position >= ? AND position < ?", i.GroceryTripID, newPosition, currPosition)
	} else {
		tx.Exec("UPDATE items SET position = position - 1 WHERE grocery_trip_id = ? AND position > ? AND position <= ?", i.GroceryTripID, currPosition, newPosition)
	}
	return nil
}

func (i *Item) AfterUpdate(tx *gorm.DB) (err error) {
	tx.Model(&GroceryTrip{}).Where("id = ?", i.GroceryTripID).Update("updated_at", time.Now())
	return nil
}
