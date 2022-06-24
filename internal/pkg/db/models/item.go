package models

import (
	"regexp"
	"strconv"
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
	StapleItemID  *uuid.UUID `gorm:"type:uuid;index"`
	Name          string     `gorm:"type:varchar(100);not null;index:idx_items_grocery_trip_id_name"`
	Quantity      int        `gorm:"default:1;not null"`
	Completed     *bool      `gorm:"default:false;not null"`
	Position      int        `gorm:"default:1;not null"`
	Notes         *string    `gorm:"type:varchar(255)"`
	MealID        *uuid.UUID `gorm:"type:uuid"`
	MealName      *string    `gorm:"type:varchar(255)"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	// Associations
	GroceryTrip GroceryTrip
	Meal        Meal
	StapleItem  StoreStapleItem
}

// BeforeCreate hook updates the item position
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

// BeforeSave hook
func (i *Item) BeforeSave(tx *gorm.DB) (err error) {
	i.Name, i.Quantity = i.parseItemName()
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

// ParseItemName handles inline quantity in the item name (e.g. Orange x 5) and
// returns a parsed version of both the name and quantity
func (i *Item) parseItemName() (parsedName string, parsedQuantity int) {
	re := regexp.MustCompile("^(.*)(\\s)x(\\s?)(\\d+)(\\s+)?")
	match := re.FindStringSubmatch(i.Name)
	if match != nil {
		var err error
		parsedQuantity, err = strconv.Atoi(match[4])
		if err != nil {
			return i.Name, i.Quantity
		}
		// Strip the quantity out of the name
		parsedName = re.ReplaceAllString(i.Name, "$1")
		return parsedName, parsedQuantity
	}
	return i.Name, i.Quantity
}
