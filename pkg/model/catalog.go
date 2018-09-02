package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// Catalog is an item in the catalog of available resources
type Catalog struct {
	gorm.Model
	Name string `gorm:"not null;unique_index"`
}

func (c Catalog) String() string {
	return fmt.Sprintf(`Catalog{Name:"%s"}`, c.Name)
}
