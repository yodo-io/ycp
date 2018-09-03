package v1

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yodo-io/ycp/pkg/model"
)

func TestGetCatalog(t *testing.T) {
	r, td := mustInitRouter(true)
	defer td()

	w := mustRequest(t, r, http.MethodGet, "/catalog")
	if w == nil {
		return
	}
	if !assert.Equal(t, http.StatusOK, w.Code) {
		return
	}

	var res []model.Catalog
	mustDecode(t, w, &res)
	assert.NotEmpty(t, res)
	for _, c := range res {
		assert.NotEmpty(t, c.Name)
	}
}
