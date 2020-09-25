package db

import (
	"log"
	"os"

	// Autoload env variables from .env
	_ "github.com/joho/godotenv/autoload"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/migration"
)

// DatabaseManager variable holds the gorm pointer to DB
type DatabaseManager struct {
	db *gorm.DB
}

var Manager *gorm.DB

// FetchConnection establishes a database connection
func FetchConnection() *gorm.DB {
	dsn := os.Getenv("POSTGRES_DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: SetLogger(),
	})
	if err != nil {
		log.Panic("[db] Database connection err: ", err)
	}
	return db
}

// SetLogger determines the log level based on DB_LOGMODE environment variable
func SetLogger() logger.Interface {
	var logLevel logger.LogLevel
	switch logMode := os.Getenv("DB_LOGMODE"); logMode {
	case "error":
		logLevel = 2
	case "warn":
		logLevel = 3
	case "info":
		logLevel = 4
	default:
		logLevel = 1
	}
	return logger.Default.LogMode(logLevel)
}

func Factory() {
	Manager = FetchConnection()

	// Automigrate on init
	log.Println("[db] Performing migrations...")
	if err := migration.AutoMigrateService(Manager); err != nil {
		log.Fatal("[db] Couldn't perform migrations! ", err)
	}
	log.Println("[db] Database connection initialized")
}
