package model

import "github.com/jinzhu/gorm"

var sampleUsers = []*User{
	{Email: "joe@example.org", Password: "secret", Role: RoleUser},
	{Email: "admin@example.org", Password: "secret", Role: RoleAdmin},
}

var sampleCatalog = []*Catalog{
	{Name: "pot.instance.small"},  // 0
	{Name: "pot.instance.large"},  // 1
	{Name: "pot.instance.xlarge"}, // 2
	{Name: "pan.instance.wok"},    // 3
	{Name: "pan.instance.s"},      // 4
	{Name: "pan.instance.m"},      // 5
	{Name: "pan.instance.xl"},     // 6
}

// UserID and CatalogID refer to idx in sample slices, will be replaced before insert
var sampleResources = []Resource{
	{Name: "pasta pot", UserID: 0, Type: "pot.instance.large"},
	{Name: "rice pot", UserID: 0, Type: "pot.instance.xlarge"},
	{Name: "stir fry pan", UserID: 1, Type: "pan.instance.wok"},
	{Name: "skillet for eggs", UserID: 1, Type: "pan.instance.s"},
}

func loadSampleData(db *gorm.DB) error {
	// Need to repeat this 3 times because we don't have a generic list type :(
	// Using []interface{} won't work: https://github.com/golang/go/wiki/InterfaceSlice

	for _, u := range sampleUsers {
		if err := db.Create(u).Error; err != nil {
			return err
		}
	}
	for _, c := range sampleCatalog {
		if err := db.Create(c).Error; err != nil {
			return err
		}
	}
	for _, rc := range sampleResources {
		// replace references with generated database IDs
		rc.UserID = sampleUsers[rc.UserID].ID
		// rc.CatalogID = sampleCatalog[rc.CatalogID].ID
		// insert
		if err := db.Create(&rc).Error; err != nil {
			return err
		}
	}
	return nil
}
