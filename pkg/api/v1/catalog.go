package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/yodo-io/ycp/pkg/model"
)

type catalog struct {
	db *gorm.DB
}

func (cc *catalog) list(c *gin.Context) (int, interface{}) {
	var entries []*model.Catalog
	if err := cc.db.Find(&entries).Error; err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, entries
}

func lookupCatalog(db *gorm.DB, name string) (*model.Catalog, error) {
	var c []*model.Catalog
	if err := db.Find(&c, "name = ?", name).Error; err != nil {
		return nil, err
	}
	if len(c) == 0 {
		return nil, nil
	}
	return c[0], nil
}
