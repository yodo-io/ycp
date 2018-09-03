package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanInsertQuota(t *testing.T) {
	db := MustInitTestDB(true)
	defer db.Close()

	q := NewQuota(1, "pot.instance.small", 10)
	if err := db.Create(&q).Error; err != nil {
		t.Fatal(err)
	}

	res := Quota{}
	db.Preload("Catalog").First(&res, q.ID)

	assert.NotZero(t, q.ID)
	assert.Equal(t, q.Type, res.Type)
	assert.Equal(t, q.Value, res.Value)
	assert.NotEmpty(t, res.UserID)
	assert.NotEmpty(t, res.Catalog)
}

func TestQuotaBelongsToUser(t *testing.T) {
	db := MustInitTestDB(true)
	defer db.Close()

	var q Quota
	var u User

	if err := db.First(&q).Related(&u).Error; err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, u)
	assert.NotEmpty(t, u.ID)
	assert.Equal(t, q.UserID, u.ID)
}
