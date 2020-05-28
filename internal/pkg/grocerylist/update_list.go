package grocerylist

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/mailer"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// UpdateListForUser updates a list for the given userID with the provided args
func UpdateListForUser(db *gorm.DB, userID uuid.UUID, args map[string]interface{}) (interface{}, error) {
	list := &models.List{}
	if err := db.Where("id = ? AND user_id = ?", args["listId"], userID).First(&list).Error; err != nil {
		return nil, err
	}

	oldName := list.Name
	if args["name"] != nil {
		list.Name = args["name"].(string)
	}
	if err := db.Save(&list).Error; err != nil {
		return nil, err
	}

	// Finally, send an email to the users of this list about this update (excluding the creator)
	if oldName != args["name"] {
		rows, err := db.Raw("SELECT u.email FROM list_users AS lu INNER JOIN users AS u ON lu.user_id = u.id WHERE lu.list_id = ? AND lu.creator = ? ORDER BY lu.created_at DESC", list.ID, false).Rows()
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var email string
		for rows.Next() {
			rows.Scan(&email)
			_, mailErr := mailer.SendListRenamedEmail(oldName, list.Name, email)
			if mailErr != nil {
				return nil, mailErr
			}
		}
	}

	return list, nil
}
