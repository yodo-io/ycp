package model

import (
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin/binding"
	validator "gopkg.in/go-playground/validator.v8"
)

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("userrole", validateRole)
	}
}

// Role constants
const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

// Role for implementing a simple RBAC model
type Role string

// User is a user in the system
// FIXME: password should be hashed
type User struct {
	ID        uint       `gorm:"primary_key"            json:"id"`
	Email     string     `gorm:"not null;unique_index"  json:"email"               binding:"required,email"`
	Password  string     `gorm:"not null"               json:"password,omitempty"  binding:"required"`
	Role      Role       `gorm:"not null"               json:"role,omitempty"      binding:"userrole"`
	Resources []Resource `json:",omitempty"`
}

func (u User) String() string {
	return fmt.Sprintf(`User{Email:"%s", Role:"%s"}`, u.Email, u.Role)
}

func validateRole(v *validator.Validate, ts reflect.Value, cs reflect.Value, f reflect.Value, ft reflect.Type, fk reflect.Kind, param string) bool {
	if val, ok := f.Interface().(Role); ok {
		// somewhat dirty, but OK for only 2 roles
		return val == "" ||
			val == RoleAdmin ||
			val == RoleUser
	}
	return true
}
