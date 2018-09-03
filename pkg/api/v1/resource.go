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

func (rc *resources) createForUser(c *gin.Context) (int, interface{}) {
	uid := c.Param("uid")

	// ensure user exists
	u, err := lookupUser(rc.db, uid)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, err
	}
	if u == nil {
		return http.StatusNotFound, err
	}

	var r model.Resource
	if err := c.ShouldBind(&r); err != nil {
		return http.StatusBadRequest, err
	}
	// lookup catalog
	cat, err := lookupCatalog(rc.db, r.Type)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, err
	}
	if cat == nil {
		return http.StatusBadRequest, errors.New("Invalid resource type")
	}

	// check quota
	ok, err := rc.checkQuota(u.ID, cat.Name)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, err
	}
	if !ok {
		return http.StatusBadRequest, errors.New("quota exceeded")
	}

	// create resource
	r.UserID = u.ID
	if err := rc.db.Create(&r).Error; err != nil {
		log.Println(err)
		return http.StatusInternalServerError, err
	}
	return http.StatusCreated, r
}

func (rc *resources) checkQuota(uid uint, tp string) (bool, error) {
	q, err := lookupQuota(rc.db, uid, tp)
	if err != nil {
		return false, err
	}
	if q == nil {
		return true, nil
	}

	var n int
	if err = rc.db.Model(&model.Resource{}).Where("user_id = ? and type = ?", uid, tp).Count(&n).Error; err != nil {
		return false, err
	}
	return n < q.Value, nil
}

func (rc *resources) listForUser(c *gin.Context) (int, interface{}) {
	userID := c.Param("uid")

	u, err := lookupUser(rc.db, userID)
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
