package main

import (
	"net/http"
	"os"
	"srv/api"
	"srv/database"
	"srv/models"

	"github.com/manyminds/api2go"
	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize logger
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)

	// Set log level from environment variable
	logLevel := getEnv("LOG_LEVEL", "info")
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		// Default to info level if parsing fails
		level = logrus.InfoLevel
		logrus.WithError(err).Warnf("Invalid LOG_LEVEL: %s, defaulting to info", logLevel)
	}
	logrus.SetLevel(level)
	logrus.Info("Starting document storage service")

	// Get database configuration from environment variables
	dbConfig := &database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "document_storage"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// Connect to database
	db, err := database.NewConnection(dbConfig)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to connect to database")
	}

	// Migrate database
	if err := database.MigrateDB(db); err != nil {
		logrus.WithError(err).Fatal("Failed to migrate database")
	}

	// Create API resources
	userResource := api.NewUserResource(db)
	folderResource := api.NewFolderResource(db)
	documentResource := api.NewDocumentResource(db)

	// Create API
	api := api2go.NewAPI("v1")

	// Register resources
	api.AddResource(models.User{}, userResource)
	api.AddResource(models.Folder{}, folderResource)
	api.AddResource(models.Document{}, documentResource)

	// Start server
	port := getEnv("PORT", "8080")
	logrus.WithField("port", port).Info("Starting server")
	http.ListenAndServe(":"+port, api.Handler())
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
