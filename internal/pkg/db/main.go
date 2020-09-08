package db

import (
	"log"
	"os"

	// Autoload env variables from .env
	_ "github.com/joho/godotenv/autoload"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/migration"
)

// ORM struct holds the gorm pointer to DB
type ORM struct {
	DB *gorm.DB
}

// FetchConnection establishes a database connection
func FetchConnection() *gorm.DB {
	dsn := os.Getenv("POSTGRES_DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic("[db] Database connection err: ", err)
	}
	return db
}

// Factory fetches a database connection, and runs some config (i.e. auto migrations)
func Factory() *ORM {
	db := FetchConnection()
	orm := &ORM{DB: db}

	// Automigrate on init
	log.Println("[db] Performing migrations...")
	if err := migration.AutoMigrateService(orm.DB); err != nil {
		log.Fatal("[db] Couldn't perform migrations! ", err)
	}
	log.Println("[db] Database connection initialized")

	return orm
}
