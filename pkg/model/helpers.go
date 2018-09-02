package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
	// sqlite test driver
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// MustInitTestDB Initialises a DB instance for testing only. As of now, this is an SQlite in-memory DB
// so it will loose state upon calling `DB.disconnect()`
// This function is for testing purpose only, it will panic if it encounters any errors.
func MustInitTestDB(sampleData bool) *gorm.DB {
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		panic(fmt.Sprintf("Error setting up db - %v", err))
	}
	db.LogMode(false) // stop spamming test outputs

	if err := Setup(db, sampleData); err != nil {
		panic(fmt.Sprintf("Error setting up db - %v", err))
	}
	return db
}
