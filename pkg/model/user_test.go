package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	db := initDB(false)
	defer db.Close()

	user := User{Email: "john@acme.org", Password: "t0ps3cr3t", Role: RoleUser}
	if err := db.Create(&user).Error; err != nil {
		t.Errorf("Expected to create user, but failed with %v", err)
	}

	var res User
	if err := db.First(&res, user.Model.ID).Error; err != nil {
		t.Fatal(err)
	}

	assert.NotZero(t, user.Model.ID)
	assert.Equal(t, res.Model.ID, user.Model.ID)
	assert.Equal(t, res.Email, user.Email)
	assert.Equal(t, res.Password, user.Password)
}

func TestEmailMustBeUnique(t *testing.T) {
	db := initDB(false)
	defer db.Close()

	user1 := User{Email: "john@acme.org", Password: "t0ps3cr3t", Role: RoleUser}
	user2 := User{Email: "john@acme.org", Password: "gu3st", Role: RoleAdmin}

	if err := db.Create(&user1).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&user2).Error; err == nil {
		t.Errorf("DB.Create(%v) should have errored but did not", user2)
	}
}

func TestListUsers(t *testing.T) {
	db := initDB(true)
	defer db.Close()

	var result []*User
	if err := db.Find(&result).Error; err != nil {
		t.Fatal(err)
	}

	assert.Len(t, result, len(sampleUsers))
}

func TestUserHasResources(t *testing.T) {
	db := initDB(true)
	defer db.Close()

	var u User
	var rcs []Resource

	if err := db.First(&u).Related(&rcs).Error; err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, rcs)
	assert.Len(t, rcs, 2)
}
