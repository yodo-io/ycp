package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type errorResponse struct {
	Error string `json:"error"`
}

// Setup registers all routes implemented by this module with the provided `RouterGroup`
// Accepting a `RouterGroup` reference makes it possible for client code to use the HTTP API
// implemented by this module along with other modules.
// The provided `gorm.DB` instance is passed to handlers to perform database related operations.
func Setup(rg *gin.RouterGroup, db *gorm.DB) {
	// user api
	uc := &users{db}
	rg.GET("/users", h(uc.list))
	rg.GET("/users/:id", h(uc.get))
	rg.POST("/users", h(uc.create))
	rg.PATCH("/users/:id", h(uc.update))
	rg.DELETE("/users/:id", h(uc.delete))

	// resource api
	rc := &resources{db}
	rg.GET("/resources/:uid", h(rc.listForUser))
	rg.GET("/resources/:uid/:rid", h(rc.getForUser))
	rg.POST("/resources/:uid", h(rc.createForUser))
	rg.PATCH("/resources/:uid/:rid", h(notImplemented)) // TODO
	rg.DELETE("/resources/:uid/:rid", h(rc.deleteForUser))

	// catalog api - can only browse for now
	cc := &catalog{db}
	rg.GET("/catalog", h(cc.list))

	// quota api
	qc := &quotas{db}
	rg.GET("/quotas/:uid", h(qc.listForUser))
	rg.POST("/quotas/:uid", h(qc.createForUser))
	rg.PATCH("/quotas/:uid/:qid", h(qc.updateForUser))
	rg.DELETE("/quotas/:uid/:qid", h(qc.deleteForUser))
}

// Simplified handler func for pure JSON APIs
// If second return val is an error it will be converted into an errorResponse
// Otherwise it will be marshalled as-is and sent along with the status code
type handlerFunc func(c *gin.Context) (int, interface{})

func notImplemented(c *gin.Context) (int, interface{}) {
	return http.StatusNotImplemented, errors.New("Not implemented")
}

// Convert internal handler funcs into gin handlers
func h(fn handlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		code, data := fn(c)
		if err, ok := data.(error); ok {
			c.JSON(code, errorResponse{Error: err.Error()})
		} else {
			c.JSON(code, data)
		}
	}
}
