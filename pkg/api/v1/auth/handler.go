package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// Handler returns the route handler for the auth module with given RouterGroup
func Handler(db *gorm.DB, secret []byte) gin.HandlerFunc {
	ac := NewController(db, secret)
	return ac.createToken
}
