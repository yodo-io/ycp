package v1

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/yodo-io/ycp/pkg/model"
)

type quotas struct {
	db *gorm.DB
}

// we can only update the value, so binding a Quota would fail for PATCH
type quotaPatch struct {
	Value int `json:"value"`
}

func (qc *quotas) listForUser(c *gin.Context) (int, interface{}) {
	var qs []*model.Quota
	uid := c.Param("uid")

	// ensure user exists
	u, err := lookupUser(qc.db, uid)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, err
	}
	if u == nil {
		return http.StatusNotFound, err
	}

	// find quotas for user
	if err := qc.db.Find(&qs, "user_id = ?", uid).Error; err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, qs
}

func (qc *quotas) createForUser(c *gin.Context) (int, interface{}) {
	uid := c.Param("uid")

	// ensure user exists
	u, err := lookupUser(qc.db, uid)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, err
	}
	if u == nil {
		return http.StatusNotFound, err
	}

	var q model.Quota
	if err := c.ShouldBind(&q); err != nil {
		return http.StatusBadRequest, err
	}

	// ensure catalog item exists
	cat, err := lookupCatalog(qc.db, q.Type)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, err
	}
	if cat == nil {
		return http.StatusBadRequest, errors.New("Invalid resource type")
	}

	// set userid and insert
	q.UserID = u.ID
	if err := qc.db.Create(&q).Error; err != nil {
		log.Println(err)
		return http.StatusInternalServerError, err
	}
	return http.StatusCreated, q
}

func (qc *quotas) updateForUser(c *gin.Context) (int, interface{}) {
	uid := c.Param("uid")
	qid := c.Param("qid")

	var qp quotaPatch
	var qs []*model.Quota
	if err := c.ShouldBind(&qp); err != nil {
		return http.StatusBadRequest, err
	}
	// if either id is wrong or id and userid don't match this will return []
	if err := qc.db.Find(&qs, "user_id = ? and id = ?", uid, qid).Error; err != nil {
		log.Println(err)
		return http.StatusInternalServerError, err
	}
	if len(qs) == 0 {
		return http.StatusNotFound, errors.New("Quota not found")
	}

	// update, prevent accidental id change
	up := gin.H{
		"value": qp.Value,
	}
	if err := qc.db.Model(&model.Quota{}).Where("id = ?", qid).Omit("id").Updates(up).Error; err != nil {
		return http.StatusInternalServerError, err
	}
	// find updated record and return
	if err := qc.db.Find(&qs, "id = ?", qid).Error; err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, qs[0]
}

func (qc *quotas) deleteForUser(c *gin.Context) (int, interface{}) {
	uid := c.Param("uid")
	qid := c.Param("qid")

	// lookup, make sure qid/uid are correct
	var qs []*model.Quota
	if err := qc.db.Find(&qs, "id = ? and user_id = ?", qid, uid).Error; err != nil {
		log.Println(err)
		return http.StatusInternalServerError, err
	}
	if len(qs) == 0 {
		return http.StatusNotFound, errors.New("Quota not found")
	}

	// delete quota
	if err := qc.db.Delete(qs[0], "id = ?", qid).Error; err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, qs[0]
}
