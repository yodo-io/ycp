package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/yodo-io/ycp/pkg/api"
	"github.com/yodo-io/ycp/pkg/api/v1"
	"github.com/yodo-io/ycp/pkg/api/v1/auth"
	"github.com/yodo-io/ycp/pkg/api/v1/rbac"
	"github.com/yodo-io/ycp/pkg/model"
)

var sampleData = true
var dbDriver = "sqlite3"
var dbString = ":memory:"
var addr = ":9000"

func main() {
	db, err := setupDB()
	if err != nil {
		log.Fatal(err)
	}
	g, err := setupGin(db)
	if err != nil {
		log.Fatal(err)
	}
	g.Run(addr)
}

func setupDB() (*gorm.DB, error) {
	db, err := gorm.Open(dbDriver, dbString)
	if err != nil {
		return nil, err
	}
	if err := model.Setup(db, sampleData); err != nil {
		return nil, err
	}
	return db, nil
}

func setupGin(db *gorm.DB) (*gin.Engine, error) {
	secret := []byte("secret")

	g := gin.Default()
	g.NoRoute(api.NotFound)

	g.POST("/auth/token", auth.Handler(db, secret))

	rg := g.Group("/v1")
	rg.Use(auth.Middleware(secret))
	rg.Use(rbac.Middleware())
	v1.Routes(rg, db)

	return g, nil
}
