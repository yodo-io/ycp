package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/yodo-io/ycp/pkg/api/test"
	"github.com/yodo-io/ycp/pkg/model"
)

// Init router for HTTP tests. Calls v1.Setup() to make sure all routes are registered
// In case of any errors this will panic, so it's not intended for use outside of test code
// Clients must defer the returned teardown function to execute any shutdown functionality
func mustInitRouter(sampleData bool) (*gin.Engine, func()) {
	db := model.MustInitTestDB(sampleData)
	teardown := func() {
		db.Close()
	}
	r := test.NewRouter()
	Setup(&r.RouterGroup, db)
	return r, teardown
}
