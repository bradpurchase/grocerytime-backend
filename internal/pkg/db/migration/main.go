package migration

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
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
		{
			// Add completed to items
			ID: "202004290846_add_completed_to_items",
			Migrate: func(tx *gorm.DB) error {
				type Item struct {
					Completed bool
				}
				return tx.AutoMigrate(&models.Item{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Table("items").DropColumn("completed").Error
			},
		},
		{
			// Add sort_order to items
			ID: "202005301245_add_position_to_items",
			Migrate: func(tx *gorm.DB) error {
				type Item struct {
					Position int
				}
				return tx.AutoMigrate(&models.Item{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Table("items").DropColumn("position").Error
			},
		},
		{
			// Add index idx_items_list_id
			ID: "202006021110_add_idx_items_list_id",
			Migrate: func(tx *gorm.DB) error {
				return tx.Table("items").AddIndex("idx_items_list_id", "list_id").Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Table("items").RemoveIndex("idx_items_list_id").Error
			},
		},
		{
			// Add index idx_list_users_list_id
			ID: "202006021117_add_idx_items_list_id",
			Migrate: func(tx *gorm.DB) error {
				return tx.Table("list_users").AddIndex("idx_list_users_list_id", "list_id").Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Table("list_users").RemoveIndex("idx_list_users_list_id").Error
			},
		},
		{
			// Add index idx_list_users_user_id
			ID: "202006021118_add_idx_items_user_id",
			Migrate: func(tx *gorm.DB) error {
				return tx.Table("list_users").AddIndex("idx_list_users_user_id", "user_id").Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Table("list_users").RemoveIndex("idx_list_users_user_id").Error
			},
		},
		{
			// Add list_id to grocery_trips
			ID: "202006040726_add_list_id_to_grocery_trips",
			Migrate: func(tx *gorm.DB) error {
				type GroceryTrip struct {
					ListID uuid.UUID
				}
				return tx.AutoMigrate(&models.GroceryTrip{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Table("grocery_trips").DropColumn("list_id").Error
			},
		},
		{
			// Add completed to grocery_trips
			ID: "202006040732_add_completed_to_grocery_trips",
			Migrate: func(tx *gorm.DB) error {
				type GroceryTrip struct {
					Completed bool
				}
				return tx.AutoMigrate(&models.GroceryTrip{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Table("grocery_trips").DropColumn("completed").Error
			},
		},
		{
			// Add grocery_trip_id to items
			ID: "202006040733_add_grocery_trip_id_to_items",
			Migrate: func(tx *gorm.DB) error {
				type Item struct {
					GroceryTripID uuid.UUID
				}
				return tx.AutoMigrate(&models.Item{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Table("items").DropColumn("grocery_trip_id").Error
			},
		},
		{
			// Add copy_remaining_items to grocery_trips
			ID: "202006071147_add_copy_remaining_items_to_grocery_trips",
			Migrate: func(tx *gorm.DB) error {
				type GroceryTrip struct {
					CopyRemainingItems bool
				}
				return tx.AutoMigrate(&models.GroceryTrip{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Table("grocery_trips").DropColumn("copy_remaining_items").Error
			},
		},
		{
			// Remove list_id from items
			ID: "202006131150_remove_list_id_from_items",
			Migrate: func(tx *gorm.DB) error {
				return tx.Table("items").DropColumn("list_id").Error
			},
			Rollback: func(tx *gorm.DB) error {
				type Item struct {
					ListID uuid.UUID
				}
				return tx.AutoMigrate(&models.Item{}).Error
			},
		},
		{
			// Add email to list_users
			ID: "202006131236_remove_list_id_from_items",
			Migrate: func(tx *gorm.DB) error {
				type ListUser struct {
					Email string
				}
				return tx.AutoMigrate(&models.ListUser{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Table("list_users").DropColumn("email").Error
			},
		},
	})
	return m.Migrate()
}
