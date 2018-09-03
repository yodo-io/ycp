package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Error string `json:"error"`
}

func Error(e error) errorResponse {
	return errorResponse{e.Error()}
}

func ErrStr(s string) errorResponse {
	return errorResponse{s}
}

func Unauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, ErrStr("Unauthorized"))
}

func Fatal(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, err)
}
