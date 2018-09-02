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
