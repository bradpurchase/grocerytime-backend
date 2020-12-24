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
			// Create default API client
			ID: "202003021034_default_api_client",
			Migrate: func(tx *gorm.DB) error {
				return tx.Create(&models.ApiClient{Name: "GroceryTime for iOS"}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Where("name = ?", "GroceryTime for iOS").Delete(&models.ApiClient{}).Error
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
			// Add password reset fields to user
			ID: "202010310840_add_password_reset_token_and_password_reset_token_expiry_to_users",
			Migrate: func(tx *gorm.DB) error {
				type User struct {
					PasswordResetToken       uuid.UUID `gorm:"type:uuid"`
					PasswordResetTokenExpiry time.Time
				}
				return tx.AutoMigrate(&User{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Exec("ALTER TABLE users DROP COLUMN password_reset_token, DROP COLUMN password_reset_token_expiry").Error
			},
		},
		{
			// Create web API client
			ID: "202010310843_web_api_client",
			Migrate: func(tx *gorm.DB) error {
				return tx.Create(&models.ApiClient{Name: "GroceryTime for Web"}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Where("name = ?", "GroceryTime for Web").Delete(&models.ApiClient{}).Error
			},
		},
		{
			// Create schema for meal planning
			ID: "202011112039_create_meal_planning_schema",
			Migrate: func(tx *gorm.DB) error {
				type Recipe struct {
					ID       uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
					UserID   uuid.UUID `gorm:"type:uuid;not null;index:idx_recipes_user_id"`
					Name     string    `gorm:"type:varchar(255);not null;index:idx_recipes_name"`
					URL      *string   `gorm:"type:varchar(255)"`
					MealType string    `gorm:"type:varchar(10);not null"`

					CreatedAt time.Time
					UpdatedAt time.Time
					DeletedAt gorm.DeletedAt
				}
				tx.AutoMigrate(&Recipe{})

				type RecipeIngredient struct {
					ID       uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
					RecipeID uuid.UUID `gorm:"type:uuid;not null;index:idx_recipe_ingredients_recipe_id"`
					Name     string    `gorm:"type:varchar(255);not null"`
					Amount   *float64  `gorm:"default:1"`
					Unit     *string   `gorm:"type:varchar(20)"`
					Notes    *string   `gorm:"type:varchar(255);"`

					CreatedAt time.Time
					UpdatedAt time.Time
					DeletedAt gorm.DeletedAt
				}
				tx.AutoMigrate(&RecipeIngredient{})

				type MealUser struct {
					ID     uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
					MealID uuid.UUID `gorm:"type:uuid;not null"`
					UserID uuid.UUID `gorm:"type:uuid;not null"`

					CreatedAt time.Time
					UpdatedAt time.Time
					DeletedAt gorm.DeletedAt
				}
				tx.AutoMigrate(&MealUser{})

				type Meal struct {
					ID       uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
					RecipeID uuid.UUID `gorm:"type:uuid;not null"`
					UserID   uuid.UUID `gorm:"type:uuid;not null"`
					Name     string    `gorm:"type:varchar(255);not null;index:idx_meals_name"`
					MealType string    `gorm:"type:varchar(10);not null"`
					Notes    *string   `gorm:"type:text"`
					Date     time.Time

					CreatedAt time.Time
					UpdatedAt time.Time
					DeletedAt gorm.DeletedAt
				}
				tx.AutoMigrate(&Meal{})

				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Exec("DROP TABLE recipes, recipe_ingredients, meals, meal_users").Error
			},
		},
		{
			// Create store_user_preferences
			ID: "202010071813_create_store_user_preferences",
			Migrate: func(tx *gorm.DB) error {
				type StoreUserPreference struct {
					ID            uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
					StoreUserID   uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"`
					DefaultStore  bool      `gorm:"default:false;not null"`
					Notifications bool      `gorm:"default:true;not null"`

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
				return tx.Exec("DELETE FROM store_user_preferences").Error
			},
		},
		{
			ID: "202011260730_add_device_name_to_auth_tokens",
			Migrate: func(tx *gorm.DB) error {
				type AuthToken struct {
					DeviceName string `gorm:"type:varchar(100)"`
				}
				return tx.AutoMigrate(&AuthToken{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Exec("ALTER TABLE auth_tokens DROP COLUMN device_name").Error
			},
		},
		{
			// Create devices
			ID: "202011291125_create_devices",
			Migrate: func(tx *gorm.DB) error {
				type Device struct {
					ID     uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
					UserID uuid.UUID `gorm:"type:uuid;not null"`
					Token  string    `gorm:"type:varchar(255);not null"`

					CreatedAt time.Time
					UpdatedAt time.Time
					DeletedAt gorm.DeletedAt
				}
				return tx.AutoMigrate(&Device{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("devices")
			},
		},
	})
	return m.Migrate()
}
