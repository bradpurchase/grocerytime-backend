package migration

import (
	"fmt"
	"log"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
)

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.ApiClient{},
		&models.AuthToken{},
		&models.Device{},
		&models.GroceryTrip{},
		&models.GroceryTripCategory{},
		&models.Item{},
		&models.Meal{},
		&models.MealUser{},
		&models.Recipe{},
		&models.RecipeIngredient{},
		&models.Store{},
		&models.StoreUser{},
		&models.StoreUserPreference{},
		&models.StoreCategory{},
		&models.User{},
	)
}

// AutoMigrateService migrates all tables and database modifications
func AutoMigrateService(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, nil)
	m.InitSchema(func(db *gorm.DB) error {
		log.Println("[Migration.InitSchema] Initializing database schema...")

		// Create extensions
		// postgres user needs to have superuser perms for now
		db.Exec("CREATE EXTENSION \"pgcrypto\";")

		if err := migrate(db); err != nil {
			return fmt.Errorf("[Migration.InitSchema]: %v", err)
		}
		return nil
	})
	m.Migrate()

	m = gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			// Create default clients
			ID: "202003021034_default_api_client",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.Create(&models.ApiClient{Name: "GroceryTime for iOS"}).Error; err != nil {
					return err
				}
				if err := tx.Create(&models.ApiClient{Name: "GroceryTime for Web"}).Error; err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Exec("TRUNCATE TABLE api_clients").Error
			},
		},
	})
	return m.Migrate()
}
