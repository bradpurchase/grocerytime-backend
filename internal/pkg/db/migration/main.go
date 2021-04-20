package migration

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/utils"
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
		{
			// Add index idx_items_grocery_trip_id_category_id
			ID: "202006021110_add_idx_items_grocery_trip_id_category_id",
			Migrate: func(tx *gorm.DB) error {
				return tx.Exec("CREATE INDEX IF NOT EXISTS idx_grocery_trip_id_category_id ON items (grocery_trip_id, category_id)").Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Exec("DROP INDEX idx_grocery_trip_id_category_id").Error
			},
		},
		{
			// Add index idx_store_users_store_id
			ID: "202006021117_add_idx_store_users_store_id_user_id",
			Migrate: func(tx *gorm.DB) error {
				return tx.Exec("CREATE INDEX IF NOT EXISTS idx_store_users_store_id_user_id ON store_users (store_id, user_id)").Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Exec("DROP INDEX idx_store_users_store_id_user_id").Error
			},
		},
		{
			// Add stores.share_code column
			ID: "202104201930_add_share_code_to_stores",
			Migrate: func(tx *gorm.DB) error {
				type Store struct {
					ShareCode string `gorm:"type:varchar(255);uniqueIndex"`
				}
				return tx.AutoMigrate(&Store{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Exec("ALTER TABLE stores DROP COLUMN share_code").Error
			},
		},
		{
			// Backfill stores.share_code column
			ID: "202104201933_backfill_stores_share_code",
			Migrate: func(tx *gorm.DB) error {
				var stores []models.Store
				if err := tx.Find(&stores).Error; err != nil {
					return err
				}
				for i := range stores {
					share_code := strings.ToUpper(utils.RandString(6))
					if err := tx.Model(&models.Store{}).Where("id = ?", stores[i].ID).Update("share_code", share_code).Error; err != nil {
						return err
					}
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Model(&models.Store{}).Not("share_code", nil).Update("share_code", nil).Error
			},
		},
	})
	return m.Migrate()
}
