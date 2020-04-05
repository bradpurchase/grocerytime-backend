package migration

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
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
			// Create default API client
			ID: "202003021034_default_api_client",
			Migrate: func(tx *gorm.DB) error {
				return tx.Create(&models.ApiClient{Name: "GroceryTime for iOS"}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Where("name = ?", "Test").Delete(&models.ApiClient{}).Error
			},
		},
	})
	return m.Migrate()
}
