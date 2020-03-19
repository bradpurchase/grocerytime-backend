package migration

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"

	uuid "github.com/satori/go.uuid"
)

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.List{},
		&models.ListUser{},
		&models.GroceryTrip{},
		&models.Item{},
		&models.ApiClient{},
		&models.AuthToken{},
	).Error
}

// AutoMigrateService migrates all tables and database modifications
func AutoMigrateService(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, nil)
	m.InitSchema(func(db *gorm.DB) error {
		log.Println("[Migration.InitSchema] Initializing database schema...")
		switch db.Dialect().GetName() {
		case "postgres":
			// Create the UUID extensions
			// postgres user needs to have superuser perms for now
			db.Exec("create extension \"pgcrypto\";")
			db.Exec("create extension \"uuid-oosp\";")
		}
		if err := migrate(db); err != nil {
			return fmt.Errorf("[Migration.InitSchema]: %v", err)
		}
		return nil
	})
	m.Migrate()

	m = gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			// Create items table
			ID: "202003011229",
			Migrate: func(tx *gorm.DB) error {
				// it's a good pratice to copy the struct inside the function,
				// so side effects are prevented if the original struct changes during the time
				type Item struct {
					ID        uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
					ListID    uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
					Name      string    `gorm:"type:varchar(100);not null"`
					Quantity  int
					CreatedAt time.Time
					UpdatedAt time.Time
				}
				return tx.AutoMigrate(&Item{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable("items").Error
			},
		},
		{
			// Create api_clients table
			ID: "202003011022",
			Migrate: func(tx *gorm.DB) error {
				type ApiClient struct {
					ID        uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
					Name      string    `gorm:"type:varchar(100);unique_index;not null"`
					Key       string    `gorm:"type:varchar(100);not null"`
					Secret    string    `gorm:"type:varchar(100);not null"`
					CreatedAt time.Time
					UpdatedAt time.Time
				}
				return tx.AutoMigrate(&ApiClient{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable("api_clients").Error
			},
		},
		{
			// Create api_tokens table
			ID: "202003011024",
			Migrate: func(tx *gorm.DB) error {
				type AuthToken struct {
					ID           uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
					ClientID     uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();not null"`
					UserID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();not null"`
					AccessToken  string    `gorm:"type:varchar(100);not null"`
					RefreshToken string    `gorm:"type:varchar(100);not null"`
					ExpiresIn    time.Time `gorm:"not null"`
					CreatedAt    time.Time
					UpdatedAt    time.Time
				}
				return tx.AutoMigrate(&AuthToken{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable("api_tokens").Error
			},
		},
		{
			// Create default API client
			ID: "202003021034",
			Migrate: func(tx *gorm.DB) error {
				return tx.Create(&models.ApiClient{Name: "Test"}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Where("name = ?", "Test").Delete(&models.ApiClient{}).Error
			},
		},
	})
	return m.Migrate()
}
