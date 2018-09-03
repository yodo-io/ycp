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
