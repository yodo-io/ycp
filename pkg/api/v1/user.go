package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/yodo-io/ycp/pkg/model"
)

type users struct {
	db *gorm.DB
}

func (uc *users) list(c *gin.Context) (int, interface{}) {
	var users []model.User
	if err := uc.db.Find(&users).Error; err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, users
}

func (uc *users) create(c *gin.Context) (int, interface{}) {
	var u model.User
	if err := c.ShouldBind(&u); err != nil {
		return http.StatusBadRequest, err
	}
	if u.Role == "" {
		u.Role = "user"
	}
	if err := uc.db.Create(&u).Error; err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusCreated, safe(u)
}

func safe(u model.User) model.User {
	u.Password = ""
	return u
}
