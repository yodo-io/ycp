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

func Fatal(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, err)
}
