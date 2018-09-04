package rbac

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/yodo-io/ycp/pkg/api/test"
	"github.com/yodo-io/ycp/pkg/api/v1/auth"
	"github.com/yodo-io/ycp/pkg/model"
)

var secret = []byte("secret")

func dummy(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func TestRBAC(t *testing.T) {
	// init db,router
	db := model.MustInitTestDB(true)
	defer db.Close()
	ac := auth.NewController(db, secret)
	r := test.NewRouter()

	// middleware for token and rbac
	rg := r.Group("/v1")
	{
		rg.Use(auth.Middleware(secret))
		rg.Use(Middleware()) // RBAC middleware

		rg.GET("/users", dummy)
		rg.GET("/users/:id", dummy)
		rg.POST("/users", dummy)
		rg.PATCH("/users/:id", dummy)
		rg.DELETE("/users/:id", dummy)

		rg.GET("/resources/:uid", dummy)
		rg.GET("/resources/:uid/:id", dummy)
		rg.POST("/resources/:uid", dummy)
		rg.PATCH("/resources/:uid/:id", dummy)
		rg.DELETE("/resources/:uid/:id", dummy)

		rg.GET("/quotas/:uid", dummy)
		rg.POST("/quotas/:uid", dummy)
		rg.DELETE("/quotas/:uid", dummy)

		rg.GET("/catalog", dummy)
	}

	tests := []struct {
		userID uint
		method string
		path   string
		code   int
	}{
		// users, user:1 - OK
		{userID: 1, method: http.MethodGet, path: "/v1/users/1", code: http.StatusOK},
		{userID: 1, method: http.MethodPatch, path: "/v1/users/1", code: http.StatusOK},
		{userID: 1, method: http.MethodDelete, path: "/v1/users/1", code: http.StatusOK},
		// // users, user:1 - Nope
		{userID: 1, method: http.MethodGet, path: "/v1/users", code: http.StatusForbidden},
		{userID: 1, method: http.MethodPost, path: "/v1/users", code: http.StatusForbidden},
		{userID: 1, method: http.MethodGet, path: "/v1/users/2", code: http.StatusForbidden},
		{userID: 1, method: http.MethodDelete, path: "/v1/users/2", code: http.StatusForbidden},
		// // catalog, user:1 - OK
		{userID: 1, method: http.MethodGet, path: "/v1/catalog", code: http.StatusOK},
		// // quotas, user:1 - OK
		{userID: 1, method: http.MethodGet, path: "/v1/quotas/1", code: http.StatusOK},
		// // quotas, user:1 - Nope
		{userID: 1, method: http.MethodPost, path: "/v1/quotas/1", code: http.StatusForbidden},
		{userID: 1, method: http.MethodDelete, path: "/v1/quotas/1", code: http.StatusForbidden},
		{userID: 1, method: http.MethodGet, path: "/v1/quotas/2", code: http.StatusForbidden},
		// // resources, user:1 - OK
		{userID: 1, method: http.MethodGet, path: "/v1/resources/1", code: http.StatusOK},
		{userID: 1, method: http.MethodGet, path: "/v1/resources/1/1", code: http.StatusOK},
		{userID: 1, method: http.MethodPatch, path: "/v1/resources/1/1", code: http.StatusOK},
		{userID: 1, method: http.MethodPost, path: "/v1/resources/1", code: http.StatusOK},
		{userID: 1, method: http.MethodGet, path: "/v1/resources/1/1", code: http.StatusOK},
		// // resources, user:1 - Nope
		{userID: 1, method: http.MethodGet, path: "/v1/resources/2", code: http.StatusForbidden},
		// as admin
		// user
		{userID: 2, method: http.MethodGet, path: "/v1/users", code: http.StatusOK},
		{userID: 2, method: http.MethodGet, path: "/v1/users/1", code: http.StatusOK},
		{userID: 2, method: http.MethodPost, path: "/v1/users", code: http.StatusOK},
		{userID: 2, method: http.MethodPatch, path: "/v1/users/1", code: http.StatusOK},
		{userID: 2, method: http.MethodDelete, path: "/v1/users/1", code: http.StatusOK},
		// resource
		{userID: 2, method: http.MethodGet, path: "/v1/resources/1", code: http.StatusOK},
		{userID: 2, method: http.MethodGet, path: "/v1/resources/2/1", code: http.StatusOK},
		{userID: 2, method: http.MethodPatch, path: "/v1/resources/2/1", code: http.StatusOK},
		{userID: 2, method: http.MethodPost, path: "/v1/resources/1", code: http.StatusOK},
		{userID: 2, method: http.MethodGet, path: "/v1/resources/1/1", code: http.StatusOK},
	}

	for _, tt := range tests {
		var u model.User
		err := db.First(&u, "id = ?", tt.userID).Error
		checkError(t, err)

		t.Log(tt.userID, tt.method, tt.path)

		tokenStr, err := ac.TokenFor(u.Email, u.Password)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(tt.method, tt.path, nil)
		req.Header.Add("Token", tokenStr)

		r.ServeHTTP(w, req)
		assert.Equal(t, tt.code, w.Code)
	}
}

func checkError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
