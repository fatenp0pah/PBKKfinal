package models

import (
	"database/sql"
	"log"
)

var db *sql.DB // Declare db as a global variable

// SetDatabase initializes the global database connection
func SetDatabase(database *sql.DB) {
	db = database
}

// GetDatabase returns the global database connection
func GetDatabase() *sql.DB {
	return db
}

// InitializeDatabase initializes the database connection and sets it globally
func InitializeDatabase() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/student_enrollment_system")
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	// Check if the connection is alive
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database: ", err)
	}
}
