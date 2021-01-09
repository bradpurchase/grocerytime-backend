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

		// Create the UUID extensions
		// postgres user needs to have superuser perms for now
		db.Exec("CREATE EXTENSION \"pgcrypto\";")
		db.Exec("CREATE EXTENSION \"uuid-oosp\";")

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
				return tx.Exec("CREATE INDEX idx_grocery_trip_id_category_id ON items (grocery_trip_id, category_id)").Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Exec("DROP INDEX idx_grocery_trip_id_category_id").Error
			},
		},
		{
			// Add index idx_store_users_store_id
			ID: "202006021117_add_idx_store_users_store_id_user_id",
			Migrate: func(tx *gorm.DB) error {
				return tx.Exec("CREATE INDEX idx_store_users_store_id_user_id ON store_users (store_id, user_id)").Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Exec("DROP INDEX idx_store_users_store_id_user_id").Error
			},
		},
		{
			// Add index idx_items_grocery_trip_id_name
			ID: "202010111143_add_idx_items_grocery_trip_id_name",
			Migrate: func(tx *gorm.DB) error {
				return tx.Exec("CREATE INDEX idx_items_name_grocery_trip_id ON items (name, grocery_trip_id)").Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Exec("DROP INDEX idx_items_name_grocery_trip_id").Error
			},
		},
		{
			// Add index for meals.date (string column)
			ID: "202101091236_meals_date_index",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.Exec("CREATE INDEX idx_meals_date ON meals (date)").Error; err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				if err := tx.Exec("DROP INDEX idx_meals_date").Error; err != nil {
					return err
				}
				return nil
			},
		},
		{
			// Add index idx_auth_tokens_access_token
			ID: "202101091252_add_idx_auth_tokens_access_token",
			Migrate: func(tx *gorm.DB) error {
				return tx.Exec("CREATE INDEX idx_auth_tokens_access_token ON auth_tokens (access_token)").Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Exec("DROP INDEX idx_auth_tokens_access_token").Error
			},
		},
	})
	return m.Migrate()
}
