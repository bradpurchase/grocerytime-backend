package user

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
)

// DeleteAccount deletes a user account
//
// This is a permanent delete, not a soft delete.
// The User model has a BeforeDelete hook to remove/clean associated data
func DeleteAccount(user models.User) (deletedUser models.User, err error) {
	if err := db.Manager.Unscoped().Delete(&user).Error; err != nil {
		return deletedUser, err
	}
	return user, err
}
