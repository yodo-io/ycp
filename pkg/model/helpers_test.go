package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func initDB(withSampleData bool) *gorm.DB {
	// Using panic here to keep usage simple and no point to recover in tests
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		panic(fmt.Sprintf("Error setting up db - %v", err))
	}

	db.LogMode(false) // stop spamming test outputs

	if err := migrateAll(db); err != nil {
		panic(fmt.Sprintf("Error setting up db - %v", err))
	}
	if withSampleData {
		if err := loadSampleData(db); err != nil {
			panic(fmt.Sprintf("Error loading sample data - %v", err))
		}
	}

	return db
}
