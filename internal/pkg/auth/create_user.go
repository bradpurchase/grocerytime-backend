package auth

import (
	"time"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser creates a user account with an email and password
func CreateUser(db *gorm.DB, email string, password string, clientID uuid.UUID) (*models.User, error) {
	passhash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	//TODO handle dupe email
	user := &models.User{
		Email:      email,
		Password:   string(passhash),
		LastSeenAt: time.Now(),
		Tokens: []models.AuthToken{
			{ClientID: clientID},
		},
		Lists: []models.List{
			{Name: "My Grocery List"},
		},
	}
	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}

	// If the user was added to a list before they signed up, we can now associate
	// those ListUser records by UserID
	db.Model(&models.ListUser{}).Where("email = ? AND user_id IS NULL").Update("user_id", user.ID)

	return user, nil
}
