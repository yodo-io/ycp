package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// Role constants
const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

// Role for implementing a simple RBAC model
type Role string

// User is a user in the system
type User struct {
	gorm.Model
	Email     string `gorm:"not null;unique_index"`
	Password  string `gorm:"not null"` // FIXME: should be hashed
	Role      Role   `gorm:"not null"`
	Resources []Resource
}

func (u User) String() string {
	return fmt.Sprintf(`User{Email:"%s", Role:"%s"}`, u.Email, u.Role)
}
