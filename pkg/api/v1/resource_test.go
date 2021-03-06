package v1

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yodo-io/ycp/pkg/api/test"
	"github.com/yodo-io/ycp/pkg/model"
)

func TestCreateResource(t *testing.T) {
	tests := []struct {
		userID uint
		in     model.Resource
		code   int
	}{
		{
			userID: 1,
			in:     model.Resource{Name: "my database", Type: "pot.instance.small"},
			code:   http.StatusCreated,
		},
		{
			userID: 10,
			in:     model.Resource{Name: "my database", Type: "pot.instance.small"},
			code:   http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		func() {
			r, td := mustInitRouter(true)
			defer td()

			w := test.MustRecord(t, r, http.MethodPost, fmt.Sprintf("/resources/%d", tt.userID), tt.in)
			if !assert.Equal(t, tt.code, w.Code) {
				return
			}
			if w.Code != http.StatusCreated {
				return
			}

			var res model.Resource
			test.MustBind(t, w, &res)
			assert.NotZero(t, res.ID)
			assert.Equal(t, tt.userID, res.UserID)
			assert.Equal(t, tt.in.Name, res.Name)
			assert.Equal(t, tt.in.Type, res.Type)
		}()
	}
}

func TestResourceQuotaLimit(t *testing.T) {
	r, td := mustInitRouter(true)
	defer td()

	limit := 10
	userID := 1

	in := model.Resource{
		Name: "a small cooking pot",
		Type: "pot.instance.small",
	}

	for i := 0; i < limit; i++ {
		w := test.MustRecord(t, r, http.MethodPost, fmt.Sprintf("/resources/%d", userID), in)
		if !assert.Equal(t, http.StatusCreated, w.Code) {
			return // unmarshalling will fail
		}
	}

	w := test.MustRecord(t, r, http.MethodPost, fmt.Sprintf("/resources/%d", userID), in)
	if !assert.Equal(t, http.StatusBadRequest, w.Code) {
		return
	}

	var e errorResponse
	test.MustBind(t, w, &e)
	assert.Regexp(t, regexp.MustCompile("quota exceeded"), e.Error)
}

func TestResourceValidation(t *testing.T) {

	tests := []struct {
		userID uint
		in     model.Resource
		err    *regexp.Regexp
	}{
		// no type
		{
			userID: 1,
			in:     model.Resource{Name: "my database"},
			err:    regexp.MustCompile(`(?i)type`),
		},
		// no name
		{
			userID: 1,
			in:     model.Resource{Type: "pot.instance.small"},
			err:    regexp.MustCompile(`(?i)name`),
		},
		// invalid type
		{
			userID: 1,
			in:     model.Resource{Name: "my foo", Type: "foo.bar.baz"},
			err:    regexp.MustCompile(`(?i)type`),
		},
	}

	for _, tt := range tests {
		func() {
			r, td := mustInitRouter(true)
			defer td()

			w := test.MustRecord(t, r, http.MethodPost, fmt.Sprintf("/resources/%d", tt.userID), tt.in)
			if !assert.Equal(t, http.StatusBadRequest, w.Code) {
				return // unmarshalling will fail
			}

			var res map[string]interface{}
			test.MustBind(t, w, &res)
			assert.NotEmpty(t, res["error"])
			assert.Regexp(t, tt.err, res["error"].(string))
		}()
	}
}

func TestGetResources(t *testing.T) {
	r, td := mustInitRouter(true)
	defer td()

	tests := []struct {
		userID uint
		code   int
	}{
		{userID: 1, code: http.StatusOK},
		{userID: 100, code: http.StatusNotFound},
	}

	for _, tt := range tests {
		// we only support to get resources scoped to a user
		w := test.MustRecord(t, r, http.MethodGet, fmt.Sprintf("/resources/%d", tt.userID))
		if !assert.Equal(t, tt.code, w.Code) {
			continue
		}
		if tt.code != http.StatusOK {
			continue
		}

		var res []*model.Resource
		test.MustBind(t, w, &res)
		assert.NotEmpty(t, res)
		for _, r := range res {
			assert.NotEmpty(t, r.Name)
			assert.NotEmpty(t, r.Type)
			assert.Equal(t, tt.userID, r.UserID)
		}
	}
}

func TestGetResource(t *testing.T) {
	r, td := mustInitRouter(true)
	defer td()

	tests := []struct {
		userID uint
		id     uint
		code   int
	}{
		{userID: 1, id: 1, code: http.StatusOK},
		{userID: 2, id: 3, code: http.StatusOK},
		{userID: 1, id: 3, code: http.StatusNotFound},
	}

	for _, tt := range tests {
		w := test.MustRecord(t, r, http.MethodGet, fmt.Sprintf("/resources/%d/%d", tt.userID, tt.id))
		if !assert.Equal(t, tt.code, w.Code) {
			continue
		}
		if w.Code != http.StatusOK {
			continue
		}

		var rc model.Resource
		test.MustBind(t, w, &rc)
		assert.NotEmpty(t, rc)
		assert.Equal(t, tt.id, rc.ID)
		assert.Equal(t, tt.userID, rc.UserID)
		assert.NotEmpty(t, rc.Name)
		assert.NotEmpty(t, rc.Type)
	}
}

func TestDeleteResourceForUser(t *testing.T) {
	r, td := mustInitRouter(true)
	defer td()

	tests := []struct {
		id     uint
		userID uint
		code   int
	}{
		{userID: 1, id: 1, code: http.StatusOK},
		{userID: 2, id: 3, code: http.StatusOK},
		{userID: 1, id: 3, code: http.StatusNotFound},
	}

	for _, tt := range tests {
		w := test.MustRecord(t, r, http.MethodDelete, fmt.Sprintf("/resources/%d/%d", tt.userID, tt.id))
		if !assert.Equal(t, tt.code, w.Code) {
			continue
		}
		if w.Code != http.StatusOK {
			continue
		}

		var rc model.Resource
		test.MustBind(t, w, &rc)
		assert.NotEmpty(t, rc)
		assert.Equal(t, tt.id, rc.ID)
		assert.Equal(t, tt.userID, rc.UserID)
		assert.NotEmpty(t, rc.Name)
		assert.NotEmpty(t, rc.Type)

		// test if resource was really deleted
		w = test.MustRecord(t, r, http.MethodGet, fmt.Sprintf("/resources/%d/%d", tt.userID, tt.id))
		assert.Equal(t, http.StatusNotFound, w.Code)
	}
}

func TestUpdateResource(t *testing.T) {
	r, td := mustInitRouter(true)
	defer td()

	tests := []struct {
		id     uint
		userID uint
		in     model.Resource
		code   int
	}{
		{
			id:     1,
			userID: 1,
			in:     model.Resource{Name: "my little cooking pot"},
			code:   http.StatusNotImplemented, // not supported yet
		},
	}

	for _, tt := range tests {
		w := test.MustRecord(t, r, http.MethodPatch, fmt.Sprintf("/resources/%d/%d", tt.userID, tt.id), tt.in)
		if !assert.Equal(t, tt.code, w.Code) {
			continue
		}
		if w.Code != http.StatusOK {
			continue
		}

		var rc model.Resource
		test.MustBind(t, w, &rc)
		assert.NotEmpty(t, rc)
		assert.Equal(t, tt.id, rc.ID)
		assert.Equal(t, tt.userID, rc.UserID)
		assert.NotEmpty(t, rc.Name)
		assert.NotEmpty(t, rc.Type)
	}
}
