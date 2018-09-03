package auth

import (
	"net/http"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/yodo-io/ycp/pkg/api/test"
	"github.com/yodo-io/ycp/pkg/model"
)

// secret to be used during testing
var secret = []byte("be00d27d0c134cc79e473f40a1e393f0")

func testKeyFunc(token *jwt.Token) (interface{}, error) {
	return secret, nil
}

func mustInitRouter() (*gin.Engine, func()) {
	db := model.MustInitTestDB(true)
	teardown := func() {
		db.Close()
	}
	r := test.NewRouter()
	Routes(&r.RouterGroup, db, secret)
	return r, teardown
}

func TestAuth(t *testing.T) {
	r, td := mustInitRouter()
	defer td()

	tests := []struct {
		tr   tokenRequest
		code int
		role model.Role
	}{
		{
			tr:   tokenRequest{Email: "joe@example.org", Password: "secret"},
			code: http.StatusOK,
			role: model.RoleUser,
		},
		{
			tr:   tokenRequest{Email: "joe", Password: "secret"},
			code: http.StatusBadRequest,
		},
		{
			tr:   tokenRequest{Password: "secret"},
			code: http.StatusBadRequest,
		},
		{
			tr:   tokenRequest{Email: "joe@example.org"},
			code: http.StatusBadRequest,
		},
		{
			tr:   tokenRequest{Email: "joe@example.org", Password: "guest"},
			code: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		w := test.MustRecord(t, r, http.MethodPost, "/token", tt.tr)
		if w == nil {
			continue
		}
		if !assert.Equal(t, tt.code, w.Code) {
			continue
		}
		if tt.code != http.StatusOK {
			continue
		}

		var res tokenResponse
		test.MustDecode(t, w, &res)

		var c Claims
		if _, err := jwt.ParseWithClaims(res.Token, &c, testKeyFunc); err != nil {
			t.Error(err)
			continue
		}

		assert.Equal(t, model.RoleUser, c.Role)
		assert.Equal(t, tt.tr.Email, c.Email)
		assert.True(t, c.StandardClaims.ExpiresAt > time.Now().Unix())
	}

}
