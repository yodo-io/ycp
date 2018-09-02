package model

import "github.com/jinzhu/gorm"

// Resource is resourced owned by a user, based on a catalog item
type Resource struct {
	gorm.Model
	Name      string `gorm:"not null"`
	UserID    uint
	CatalogID uint
}
