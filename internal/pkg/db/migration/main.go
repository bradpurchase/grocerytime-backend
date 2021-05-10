package migration

import (
	"fmt"
	"log"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/gofrs/uuid"
	"gorm.io/datatypes"
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
			// Create store_item_category_settings
			ID: "202105020850_create_store_item_category_settings",
			Migrate: func(tx *gorm.DB) error {
				type StoreItemCategorySettings struct {
					ID      uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
					StoreID uuid.UUID `gorm:"type:uuid;not null;index:idx_store_category_items_store_id"`
					Items   datatypes.JSON

					CreatedAt time.Time
					UpdatedAt time.Time
					DeletedAt gorm.DeletedAt
				}
				return tx.AutoMigrate(&StoreItemCategorySettings{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("store_item_category_settings")
			},
		},
		{
			ID: "202105050742_create_store_staple_items",
			Migrate: func(tx *gorm.DB) error {
				type StoreStapleItem struct {
					ID      uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
					StoreID uuid.UUID `gorm:"type:uuid;not null;index:idx_store_staple_items_store_id"`
					Name    string    `gorm:"type:varchar(100);not null"`

					CreatedAt time.Time
					UpdatedAt time.Time
					DeletedAt gorm.DeletedAt
				}
				return tx.AutoMigrate(&StoreStapleItem{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("store_staple_items")
			},
		},
		{
			// Add index idx_store_staple_items_store_id_name
			ID: "202105050820_add_idx_store_staple_items_store_id_name",
			Migrate: func(tx *gorm.DB) error {
				return tx.Exec("CREATE INDEX IF NOT EXISTS idx_store_staple_items_store_id_name ON store_staple_items (store_id, name)").Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Exec("DROP INDEX idx_store_staple_items_store_id_name").Error
			},
		},
		{
			// Add column staple_id to items
			ID: "202105060735_add_staple_id_to_items",
			Migrate: func(tx *gorm.DB) error {
				type Item struct {
					StapleItemID *uuid.UUID `gorm:"type:uuid;index"`
				}
				return tx.AutoMigrate(&Item{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropColumn("people", "age")
			},
		},
	})
	return m.Migrate()
}
