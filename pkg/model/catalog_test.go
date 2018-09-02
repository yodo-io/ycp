package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanInsertCatalog(t *testing.T) {
	db := initDB(false)
	defer db.Close()

	item := Catalog{Name: "db.instance.small"}
	if err := db.Create(&item).Error; err != nil {
		t.Fatal(db.Error)
	}

	res := Catalog{}
	db.First(&res, item.Model.ID)

	assert.NotZero(t, item.Model.ID)
	assert.Equal(t, item.Name, res.Name)
}

func TestCannotInsertDuplicateName(t *testing.T) {
	db := initDB(false)
	defer db.Close()

	name := "db.instance.small"

	item1 := Catalog{Name: name}
	item2 := Catalog{Name: name}

	if err := db.Create(&item1).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&item2).Error; err == nil {
		t.Errorf("DB.Create(%v) should have errored but did not", item2)
	}
}

func TestListCatalog(t *testing.T) {
	db := initDB(true)
	defer db.Close()

	var result []*Catalog
	if err := db.Find(&result).Error; err != nil {
		t.Fatal(err)
	}

	assert.Len(t, result, len(sampleCatalog))
}
