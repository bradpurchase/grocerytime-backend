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
		db.Exec("create extension \"pgcrypto\";")
		db.Exec("create extension \"uuid-oosp\";")

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
				return tx.Migrator().CreateIndex(&models.Item{}, "idx_grocery_trip_id_category_id")
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropIndex(&models.Item{}, "idx_grocery_trip_id_category_id")
			},
		},
		{
			// Add index idx_store_users_store_id
			ID: "202006021117_add_idx_store_users_store_id_user_id",
			Migrate: func(tx *gorm.DB) error {
				return tx.Migrator().CreateIndex(&models.StoreUser{}, "idx_store_users_store_id_user_id")
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropIndex(&models.StoreUser{}, "idx_store_users_store_id_user_id")
			},
		},
		{
			// Add index idx_store_users_user_id
			ID: "202006021118_add_idx_store_users_store_id",
			Migrate: func(tx *gorm.DB) error {
				return tx.Migrator().CreateIndex(&models.StoreUser{}, "idx_store_users_store_id")
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropIndex(&models.StoreUser{}, "idx_store_users_store_id")
			},
		},
		{
			// Add index idx_store_categories_store_id
			ID: "202009060829_add_idx_store_categories_store_id",
			Migrate: func(tx *gorm.DB) error {
				return tx.Migrator().CreateIndex(&models.StoreCategory{}, "idx_store_categories_store_id")
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropIndex(&models.StoreCategory{}, "idx_store_categories_store_id")
			},
		},
		{
			// Add index idx_grocery_trips_store_id
			ID: "202009060831_add_idx_grocery_trips_store_id",
			Migrate: func(tx *gorm.DB) error {
				return tx.Migrator().CreateIndex(&models.GroceryTrip{}, "idx_grocery_trips_store_id")
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropIndex(&models.GroceryTrip{}, "idx_grocery_trips_store_id")
			},
		},
		{
			// Add index idx_grocery_trip_categories_grocery_trip_id
			ID: "202009060833_add_idx_grocery_trip_categories_grocery_trip_id",
			Migrate: func(tx *gorm.DB) error {
				return tx.Migrator().CreateIndex(&models.GroceryTripCategory{}, "idx_grocery_trip_categories_grocery_trip_id")
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropIndex(&models.GroceryTripCategory{}, "idx_grocery_trip_categories_grocery_trip_id")
			},
		},
		{
			// Add index idx_grocery_trip_categories_grocery_trip_id
			ID: "202009060835_add_idx_grocery_trip_categories_store_category_id",
			Migrate: func(tx *gorm.DB) error {
				return tx.Migrator().CreateIndex(&models.GroceryTripCategory{}, "idx_grocery_trip_categories_store_category_id")
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropIndex(&models.GroceryTripCategory{}, "idx_grocery_trip_categories_store_category_id")
			},
		},
	})
	return m.Migrate()
}
