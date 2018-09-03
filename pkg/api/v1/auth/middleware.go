package auth

import (
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/yodo-io/ycp/pkg/api"
)

// Middleware returns a gin.HandlerFunc implementing auth middleware for the application
func Middleware(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Token")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrStr("Unauthorized"))
			return
		}
		cl, err := parseToken(token, secret)
		if err != nil {
			handleTokenError(c, err)
			return
		}
		c.Set("claims", cl)
	}
}

func handleTokenError(c *gin.Context, err error) {
	if _, ok := err.(*jwt.ValidationError); ok {
		api.Unauthorized(c)
	}
	api.Fatal(c, err)
}

func parseToken(tokenString string, secret []byte) (Claims, error) {
	var claims Claims
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		// validate the algo is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	return claims, err
}
