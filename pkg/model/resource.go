package model

// Resource is resourced owned by a user, based on a catalog item
type Resource struct {
	ID      uint   `gorm:"primary_key"`
	Name    string `gorm:"not null"               binding:"required"`
	UserID  uint
	Type    string  `                             binding:"required"`
	Catalog Catalog `gorm:"foreignkey:Type"`
}
