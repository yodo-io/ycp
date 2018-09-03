package auth

import (
	"errors"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/yodo-io/ycp/pkg/api"
	"github.com/yodo-io/ycp/pkg/model"
)

var errAuthFailed = errors.New("Authentication failed")

var tokenLifetime = 15 * time.Minute
var tokenIssuer = "ycp"

type tokenRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type tokenResponse struct {
	Token string `json:"token" binding:"required"`
}

type authz struct {
	db     *gorm.DB
	secret []byte
}

// Claims contains claims attached to the auth token. They will be stored in the
// gin.Context upon successful validation of the user provided token
type Claims struct {
	Role  model.Role `json:"role"`
	Email string     `json:"email"`
	jwt.StandardClaims
}

func newAuthz(db *gorm.DB, secret []byte) *authz {
	return &authz{db, secret}
}

func newResponse(tokenStr string) *tokenResponse {
	return &tokenResponse{tokenStr}
}

func newRequest(email, password string) *tokenRequest {
	return &tokenRequest{
		Email:    email,
		Password: password,
	}
}

func claimsFor(u *model.User) *Claims {
	return &Claims{
		Role:  u.Role,
		Email: u.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenLifetime).Unix(),
			Issuer:    tokenIssuer,
		},
	}
}

func (a *authz) tokenFor(tr *tokenRequest) (string, error) {
	u, err := a.validateUser(tr.Email, tr.Password)
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsFor(u))
	return token.SignedString(a.secret)
}

func (a *authz) validateUser(email string, pw string) (*model.User, error) {
	var u []*model.User
	if err := a.db.Find(&u, "email = ?", email).Error; err != nil {
		return nil, err
	}
	if len(u) == 0 || u[0].Password != pw {
		return nil, errAuthFailed
	}
	return u[0], nil
}

func (a *authz) createToken(c *gin.Context) {
	var tr tokenRequest
	if err := c.ShouldBind(&tr); err != nil {
		c.JSON(http.StatusBadRequest, api.Error(err))
		return
	}

	ts, err := a.tokenFor(&tr)
	if err == errAuthFailed {
		c.JSON(http.StatusBadRequest, api.Error(err))
		return
	}
	if err != nil {
		api.Fatal(c, err)
		return
	}

	c.JSON(http.StatusOK, newResponse(ts))
}
