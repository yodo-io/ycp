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

// Auth implements authentication for the API
type Auth struct {
	db     *gorm.DB
	secret []byte
}

// Claims contains claims attached to the auth token. They will be stored in the
// gin.Context upon successful validation of the user provided token
type Claims struct {
	Role   model.Role `json:"role"`
	UserID uint       `json:"userID"`
	Email  string     `json:"email"`
	jwt.StandardClaims
}

// NewController create a new auth controller
func NewController(db *gorm.DB, secret []byte) *Auth {
	return &Auth{db, secret}
}

func newResponse(tokenStr string) *tokenResponse {
	return &tokenResponse{tokenStr}
}

// NewRequest generate a new tokenRequest
func newRequest(email, password string) *tokenRequest {
	return &tokenRequest{
		Email:    email,
		Password: password,
	}
}

func claimsFor(u *model.User) *Claims {
	return &Claims{
		Role:   u.Role,
		Email:  u.Email,
		UserID: u.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenLifetime).Unix(),
			Issuer:    tokenIssuer,
		},
	}
}

// TokenFor generates a new token for given tokenRequest
// Might be better to have this in a dedicated component (TokenProvider or sthg.) instead of making
// the entire controller public.
func (a *Auth) TokenFor(email, password string) (string, error) {
	u, err := a.validateUser(email, password)
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsFor(u))
	return token.SignedString(a.secret)
}

func (a *Auth) validateUser(email string, pw string) (*model.User, error) {
	var u []*model.User
	if err := a.db.Find(&u, "email = ?", email).Error; err != nil {
		return nil, err
	}
	if len(u) == 0 || u[0].Password != pw {
		return nil, errAuthFailed
	}
	return u[0], nil
}

func (a *Auth) createToken(c *gin.Context) {
	var tr tokenRequest
	if err := c.ShouldBind(&tr); err != nil {
		c.JSON(http.StatusBadRequest, api.Error(err))
		return
	}

	ts, err := a.TokenFor(tr.Email, tr.Password)
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
