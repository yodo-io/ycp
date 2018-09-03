package v1

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/yodo-io/ycp/pkg/model"
)

type resources struct {
	db *gorm.DB
}

func (rc *resources) create(c *gin.Context) (int, interface{}) {
	var r model.Resource
	if err := c.ShouldBind(&r); err != nil {
		return http.StatusBadRequest, err
	}
	// lookup catalog
	cat, err := rc.lookupCatalog(r.Type)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, err
	}
	if cat == nil {
		return http.StatusBadRequest, errors.New("Invalid resource type")
	}
	// create resource
	if err := rc.db.Create(&r).Error; err != nil {
		log.Println(err)
		return http.StatusInternalServerError, err
	}
	return http.StatusCreated, r
}

func (rc *resources) listForUser(c *gin.Context) (int, interface{}) {
	userID := c.Param("uid")

	u, err := rc.lookupUser(userID)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, err
	}
	if u == nil {
		return http.StatusNotFound, err
	}

	var rs []*model.Resource
	if err := rc.db.Model(&u).Related(&rs).Error; err != nil {
		log.Println(err)
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, rs
}

func (rc *resources) getForUser(c *gin.Context) (int, interface{}) {
	uid := c.Param("uid")
	rid := c.Param("rid")

	var rs []*model.Resource
	if err := rc.db.Find(&rs, "id = ? and user_id = ?", rid, uid).Error; err != nil {
		log.Println(err)
		return http.StatusInternalServerError, err
	}
	if len(rs) == 0 {
		return http.StatusNotFound, errors.New("Resource not found")
	}
	return http.StatusOK, rs[0]
}

func (rc *resources) deleteForUser(c *gin.Context) (int, interface{}) {
	uid := c.Param("uid")
	rid := c.Param("rid")

	// lookup, make sure rid/uid are correct
	var rs []*model.Resource
	if err := rc.db.Find(&rs, "id = ? and user_id = ?", rid, uid).Error; err != nil {
		log.Println(err)
		return http.StatusInternalServerError, err
	}
	if len(rs) == 0 {
		return http.StatusNotFound, errors.New("Resource not found")
	}

	// delete resource
	if err := rc.db.Delete(rs[0], "id = ?", rid).Error; err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, rs[0]
}

func (rc *resources) lookupUser(id string) (*model.User, error) {
	var u []*model.User
	if err := rc.db.Find(&u, "id = ?", id).Error; err != nil {
		return nil, err
	}
	if len(u) == 0 {
		return nil, nil
	}
	return u[0], nil
}

func (rc *resources) lookupCatalog(name string) (*model.Catalog, error) {
	var c []*model.Catalog
	if err := rc.db.Find(&c, "name = ?", name).Error; err != nil {
		return nil, err
	}
	if len(c) == 0 {
		return nil, nil
	}
	return c[0], nil
}
