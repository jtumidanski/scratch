package database

import (
	"fmt"
	"srv/models"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Config holds the database configuration
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewConnection creates a new database connection
func NewConnection(config *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// NewSQLiteConnection creates a new SQLite database connection for testing
func NewSQLiteConnection(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// MigrateDB performs database migration
func MigrateDB(db *gorm.DB) error {
	logrus.Info("Running database migrations")

	// Auto migrate the models
	err := db.AutoMigrate(
		&models.User{},
		&models.Folder{},
		&models.Document{},
	)

	if err != nil {
		logrus.WithError(err).Error("Failed to migrate database")
		return err
	}

	logrus.Info("Database migration completed successfully")
	return nil
}
