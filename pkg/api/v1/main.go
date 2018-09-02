package v1

import (
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
	uc := &users{db}

	rg.GET("/users", h(uc.list))
	rg.POST("/users", h(uc.create))
}

// Simplified handler func for pure JSON APIs
// If second return val is an error it will be converted into an errorResponse
// Otherwise it will be marshalled as-is and sent along with the status code
type handlerFunc func(c *gin.Context) (int, interface{})

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
