package model

import (
	"github.com/jinzhu/gorm"
)

// migrateAll migrates the given database's schema for all model types
// auto migration will only add missing fields, won't delete/change current data
func migrateAll(db *gorm.DB) error {
	return db.AutoMigrate(
		&Catalog{},
		&User{},
		&Resource{},
	).Error
}

// Setup migrates the DB and optionally loads sample data
func Setup(db *gorm.DB, sampleData bool) error {
	if err := migrateAll(db); err != nil {
		return err
	}
	if sampleData {
		if err := loadSampleData(db); err != nil {
			return err
		}
	}
	return nil
}
