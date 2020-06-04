package models

import (
	"time"

	// Postgres dialect for GORM
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	uuid "github.com/satori/go.uuid"
)

// Item defines the model for items
type Item struct {
	ID            uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	ListID        uuid.UUID `gorm:"type:uuid;not null"`
	GroceryTripID uuid.UUID `gorm:"type:uuid;not null"`
	UserID        uuid.UUID `gorm:"type:uuid;not null"`
	Name          string    `gorm:"type:varchar(100);not null"`
	Quantity      int       `gorm:"default:1;not null"`
	Completed     bool      `gorm:"default:false;not null"`
	Position      int       `gorm:"default:1;not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time

	// Associations
	List List
}

func (i *Item) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Exec("UPDATE items SET position = position + 1 WHERE list_id = ? AND position >= 0", i.ListID)
	return nil
}

// AfterCreate hook to touch the associated list after an item is created
// so that its UpdatedAt column is updated
//
// Note: We're updating an arbitrary column here to get UpdatedAt to update -
// not sure if this is needed or if there's a better way to do this...
func (i *Item) AfterCreate(tx *gorm.DB) (err error) {
	tx.Model(&List{}).Where("id = ?", i.ListID).Update("updated_at", time.Now())
	return nil
}

// BeforeUpdate hook handles reordering item around an item when
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
		tx.Exec("UPDATE items SET position = position + 1 WHERE list_id = ? AND position >= ? AND position < ?", i.ListID, newPosition, currPosition)
	} else {
		tx.Exec("UPDATE items SET position = position - 1 WHERE list_id = ? AND position > ? AND position <= ?", i.ListID, currPosition, newPosition)
	}
	return nil
}

func (i *Item) AfterUpdate(tx *gorm.DB) (err error) {
	tx.Model(&List{}).Where("id = ?", i.ListID).Update("updated_at", time.Now())
	return nil
}

func (i *Item) AfterDelete(tx *gorm.DB) (err error) {
	tx.Model(&List{}).Where("id = ?", i.ListID).Update("updated_at", time.Now())
	return nil
}
