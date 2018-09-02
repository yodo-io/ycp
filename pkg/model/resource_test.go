package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanInsertResource(t *testing.T) {
	db := MustInitTestDB(false)
	defer db.Close()

	rc := Resource{Name: "my database server"}
	if err := db.Create(&rc).Error; err != nil {
		t.Fatal(db.Error)
	}

	res := Resource{}
	db.First(&res, rc.Model.ID)

	assert.NotZero(t, rc.Model.ID)
	assert.Equal(t, rc.Name, res.Name)
}

func TestListResources(t *testing.T) {
	db := MustInitTestDB(true)
	defer db.Close()

	var result []*Resource
	if err := db.Find(&result).Error; err != nil {
		t.Fatal(err)
	}

	assert.Len(t, result, len(sampleResources))
}

func TestBelongsToUser(t *testing.T) {
	db := MustInitTestDB(true)
	defer db.Close()

	var rc Resource
	var u User

	if err := db.First(&rc).Related(&u).Error; err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, u)
	assert.Equal(t, "joe@example.org", u.Email)
}

func TestBelongsToCatalog(t *testing.T) {
	db := MustInitTestDB(true)
	defer db.Close()

	var rc Resource
	var c Catalog

	if err := db.First(&rc).Related(&c).Error; err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, c)
	assert.Equal(t, "pot.instance.large", c.Name)
}
