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
		&models.GroceryTrip{},
		&models.GroceryTripCategory{},
		&models.Item{},
		&models.Store{},
		&models.StoreUser{},
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
			// Create default API client
			ID: "202003021034_default_api_client",
			Migrate: func(tx *gorm.DB) error {
				return tx.Create(&models.ApiClient{Name: "GroceryTime for iOS"}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Where("name = ?", "Test").Delete(&models.ApiClient{}).Error
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
			// Add index idx_store_users_user_id
			ID: "202006021118_add_idx_store_users_store_id",
			Migrate: func(tx *gorm.DB) error {
				return tx.Exec("CREATE INDEX idx_store_users_store_id ON store_users (store_id)").Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Exec("DROP INDEX idx_store_users_store_id").Error
			},
		},
		{
			// Add index idx_store_categories_store_id
			ID: "202009060829_add_idx_store_categories_store_id",
			Migrate: func(tx *gorm.DB) error {
				return tx.Exec("CREATE INDEX idx_store_categories_store_id ON store_categories (store_id)").Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Exec("DROP INDEX idx_store_categories_store_id").Error
			},
		},
		{
			// Add index idx_grocery_trips_store_id
			ID: "202009060831_add_idx_grocery_trips_store_id",
			Migrate: func(tx *gorm.DB) error {
				return tx.Exec("CREATE INDEX idx_grocery_trips_store_id ON grocery_trips (store_id)").Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Exec("DROP INDEX idx_grocery_trips_store_id").Error
			},
		},
		{
			// Add index idx_grocery_trip_categories_grocery_trip_id
			ID: "202009060833_add_idx_grocery_trip_categories_grocery_trip_id",
			Migrate: func(tx *gorm.DB) error {
				return tx.Exec("CREATE INDEX idx_grocery_trip_categories_grocery_trip_id ON grocery_trip_categories (grocery_trip_id)").Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Exec("DROP INDEX idx_grocery_trip_categories_grocery_trip_id").Error
			},
		},
		{
			// Add index idx_grocery_trip_categories_grocery_trip_id
			ID: "202009060835_add_idx_grocery_trip_categories_store_category_id",
			Migrate: func(tx *gorm.DB) error {
				return tx.Exec("CREATE INDEX idx_grocery_trip_categories_store_category_id ON grocery_trip_categories (store_category_id)").Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Exec("DROP INDEX idx_grocery_trip_categories_store_category_id").Error
			},
		},
		{
			// Rename user first_name to name
			ID: "202009161947_rename_user_first_name",
			Migrate: func(tx *gorm.DB) error {
				return tx.Migrator().RenameColumn(&models.User{}, "first_name", "name")
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().RenameColumn(&models.User{}, "name", "first_name")
			},
		},
		{
			// Drop last_name from users
			ID: "202009161947_drop_last_name",
			Migrate: func(tx *gorm.DB) error {
				return tx.Migrator().DropColumn(&models.User{}, "last_name")
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().AddColumn(&models.User{}, "last_name")
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
	})
	return m.Migrate()
}
