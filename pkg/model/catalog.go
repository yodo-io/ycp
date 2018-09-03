package model

import (
	"fmt"
)

// Catalog is an item in the catalog of available resources
type Catalog struct {
	Name string `gorm:"not null;primary_key"`
}

func (c Catalog) String() string {
	return fmt.Sprintf(`Catalog{Name:"%s"}`, c.Name)
}
