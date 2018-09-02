package model

import (
	"fmt"
)

// Catalog is an item in the catalog of available resources
type Catalog struct {
	ID   uint   `gorm:"primary_key"`
	Name string `gorm:"not null;unique_index"`
}

func (c Catalog) String() string {
	return fmt.Sprintf(`Catalog{Name:"%s"}`, c.Name)
}
