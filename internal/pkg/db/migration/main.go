package migration

import (
	"fmt"
	"log"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	uuid "github.com/satori/go.uuid"
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
			// Create store_user_preferences
			ID: "202010071813_create_store_user_preferences",
			Migrate: func(tx *gorm.DB) error {
				type StoreUserPreference struct {
					ID            uuid.UUID
					StoreUserID   uuid.UUID
					DefaultStore  bool
					Notifications bool

					CreatedAt time.Time
					UpdatedAt time.Time
					DeletedAt gorm.DeletedAt
				}
				return tx.AutoMigrate(&StoreUserPreference{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("store_user_preferences")
			},
		},
		{
			// Create missing store_user_preferences records
			ID: "202010071829_create_store_user_preferences_records",
			Migrate: func(tx *gorm.DB) error {
				var storeUsers []models.StoreUser
				if err := tx.Where("active = ?", true).Find(&storeUsers).Error; err != nil {
					return err
				}
				for i := range storeUsers {
					storeUserPref := &models.StoreUserPreference{
						ID:          uuid.NewV4(),
						StoreUserID: storeUsers[i].ID,
					}
					if err := tx.Create(&storeUserPref).Error; err != nil {
						return err
					}
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				// Empty all records in the table
				var storeUserPrefs []models.StoreUserPreference
				if err := tx.Find(&storeUserPrefs).Error; err != nil {
					return err
				}
				return nil
			},
		},
	})
	return m.Migrate()
}
