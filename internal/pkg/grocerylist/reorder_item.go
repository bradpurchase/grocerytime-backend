package grocerylist

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/jinzhu/gorm"
)

// ReorderItem handles the reordering of an item in the list by taking the
// item ID and the new position. It returns he reordered list.
func ReorderItem(db *gorm.DB, itemID interface{}, position int) (*models.List, error) {
	item := &models.Item{}
	if err := db.Where("id = ?", itemID).First(&item).Error; err != nil {
		return nil, err
	}
	item.Position = position
	if err := db.Save(&item).Error; err != nil {
		return nil, err
	}
	list := &models.List{}
	if err := db.Where("id = ?", item.ListID).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
