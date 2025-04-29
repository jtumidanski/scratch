package database

import (
	"srv/models"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// NewTestDB creates a new in-memory SQLite database for testing
func NewTestDB(t *testing.T) *gorm.DB {
	db, err := NewSQLiteConnection("file::memory:?cache=shared")
	require.NoError(t, err, "Failed to connect to test database")

	// Migrate the schema
	err = db.AutoMigrate(
		&models.User{},
		&models.Folder{},
		&models.Document{},
	)
	require.NoError(t, err, "Failed to migrate test database")

	return db
}

// CleanupTestDB cleans up the test database
func CleanupTestDB(t *testing.T, db *gorm.DB) {
	// Get the underlying SQL database
	sqlDB, err := db.DB()
	require.NoError(t, err, "Failed to get SQL database")

	// Close the database connection
	err = sqlDB.Close()
	require.NoError(t, err, "Failed to close test database")
}

// TruncateTables truncates all tables in the test database
func TruncateTables(t *testing.T, db *gorm.DB) {
	require.NoError(t, db.Exec("DELETE FROM documents").Error, "Failed to truncate documents table")
	require.NoError(t, db.Exec("DELETE FROM folders").Error, "Failed to truncate folders table")
	require.NoError(t, db.Exec("DELETE FROM users").Error, "Failed to truncate users table")
}