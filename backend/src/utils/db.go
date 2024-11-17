package utils

import (
	"Delingo/src/models"
	"database/sql"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var GormDB *gorm.DB
var SQLDB *sql.DB

// Modify InitDB to return an error
func InitDB() error {
	var err error
	dsn := "user=postgres password=aine dbname=delingo1 host=localhost port=5433 sslmode=disable"

	// Initialize GORM
	GormDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err // Return error if connection fails
	}

	// Get the underlying sql.DB connection
	SQLDB, err = GormDB.DB()
	if err != nil {
		return err // Return error if getting raw DB connection fails
	}

	// Test the connection
	if err := SQLDB.Ping(); err != nil {
		return err // Return error if ping fails
	}

	// Auto-migrate GORM models
	if err := GormDB.AutoMigrate(&models.User{}, &models.Profile{}); err != nil {
		return err // Return error if migration fails
	}

	return nil // No error, successful initialization
}
