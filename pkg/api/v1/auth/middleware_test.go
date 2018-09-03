package auth

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/yodo-io/ycp/pkg/api/test"
	"github.com/yodo-io/ycp/pkg/model"
)

func mustInitMiddleware() *gin.Engine {
	r := test.NewRouter()
	r.Use(Middleware(secret))
	return r
}

func TestAuthMiddleware(t *testing.T) {
	db := model.MustInitTestDB(true)
	defer db.Close()
	ac := newAuthz(db, secret)

	tokenStr, err := ac.tokenFor(newRequest("joe@example.org", "secret"))
	if err != nil {
		t.Fatal(err)
		return
	}

	tests := []struct {
		token string
		code  int
	}{
		{
			token: tokenStr,
			code:  http.StatusOK,
		},
		{
			token: "",
			code:  http.StatusUnauthorized,
		},
		{
			token: "foobar",
			code:  http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		func() {
			var claims interface{}

			n := 0
			r := mustInitMiddleware()

			r.GET("/private", func(c *gin.Context) {
				n++
				claims, _ = c.Get("claims")
				c.JSON(http.StatusOK, gin.H{"message": "hello"})
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/private", nil)
			req.Header.Add("Token", tt.token)

			r.ServeHTTP(w, req)

			// check status code
			if tt.code != w.Code {
				t.Errorf("Expected code %d, got %d", tt.code, w.Code)
			}
			// if not authorized, should not call route handler
			if tt.code != http.StatusOK && n > 0 {
				t.Errorf("Expected handler to not be called, but was called %d times", n)
			}
			// done here if code != 200
			if tt.code != http.StatusOK {
				return
			}
			// check claims have been set
			if _, ok := claims.(Claims); !ok {
				t.Errorf("Expected claims to be of type Claims but was %v", reflect.TypeOf(claims))
			}
		}()
	}
}
