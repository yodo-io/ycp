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
	var users []*model.User
	if err := uc.db.Find(&users).Error; err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, scrubAll(users)
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
	return http.StatusCreated, scrub(&u)
}

func (uc *users) get(c *gin.Context) (int, interface{}) {
	var u []*model.User
	id := c.Param("id")

	if err := uc.db.Find(&u, "id = ?", id).Error; err != nil {
		return http.StatusInternalServerError, err
	}
	if len(u) == 0 {
		return http.StatusNotFound, errorResponse{Error: "Not Found"}
	}
	return http.StatusOK, scrub(u[0])
}

func (uc *users) delete(c *gin.Context) (int, interface{}) {
	id := c.Param("id")
	var u []*model.User

	if err := uc.db.Find(&u, "id = ?", id).Error; err != nil {
		return http.StatusInternalServerError, err
	}
	if len(u) == 0 {
		return http.StatusNotFound, errorResponse{Error: "Not Found"}
	}
	if err := uc.db.Delete(u[0], "id = ?", id).Error; err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, scrub(u[0])
}

type userPatch struct {
	Email    string     `json:"email" binding:"omitempty,email"`
	Role     model.Role `json:"role" binding:"omitempty,userrole"`
	Password string     `json:"password"`
}

func (uc *users) update(c *gin.Context) (int, interface{}) {
	id := c.Param("id")
	var u []*model.User

	if err := uc.db.Find(&u, "id = ?", id).Error; err != nil {
		return http.StatusInternalServerError, err
	}
	if len(u) == 0 {
		return http.StatusNotFound, errorResponse{"Not found"}
	}

	var up userPatch
	if err := c.ShouldBind(&up); err != nil {
		return http.StatusBadRequest, err
	}
	if err := uc.db.Model(&model.User{}).Where("id = ?", id).Omit("id").Updates(up).Error; err != nil {
		return http.StatusInternalServerError, err
	}
	if err := uc.db.Find(&u, "id = ?", id).Error; err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, scrub(u[0])
}

// Remove sensitive information from user object
func scrub(u *model.User) *model.User {
	u.Password = ""
	return u
}

// Remove sensitive information from a sclice of user objects
func scrubAll(users []*model.User) []*model.User {
	for _, u := range users {
		u.Password = ""
	}
	return users
}
