package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// Routes registers the routes for the auth module with given RouterGroup
func Routes(rg *gin.RouterGroup, db *gorm.DB, secret []byte) {
	ac := new(db, secret)
	rg.POST("/token", ac.createToken)
}
