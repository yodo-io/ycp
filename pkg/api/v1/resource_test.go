package v1

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yodo-io/ycp/pkg/model"
)

func TestCreateResource(t *testing.T) {
	tests := []struct {
		in   model.Resource
		out  model.Resource
		code int
	}{
		{
			in:   model.Resource{Name: "my database", Type: "pot.instance.small"},
			out:  model.Resource{Name: "my database", Type: "pot.instance.small", ID: 5},
			code: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		func() {
			r, td := mustInitRouter(true)
			defer td()

			w := mustRequest(t, r, http.MethodPost, "/resources", tt.in)
			if w == nil {
				return
			}
			if !assert.Equal(t, tt.code, w.Code) {
				return // unmarshalling will fail
			}

			var res model.Resource
			mustDecode(t, w, &res)
			assert.NotZero(t, res.ID)
			assert.Equal(t, tt.out.Name, res.Name)
			assert.Equal(t, tt.out.Type, res.Type)
		}()
	}
}

func TestResourceValidation(t *testing.T) {

	tests := []struct {
		in  model.Resource
		err *regexp.Regexp
	}{
		// no type
		{
			in:  model.Resource{Name: "my database"},
			err: regexp.MustCompile(`(?i)type`),
		},
		// no name
		{
			in:  model.Resource{Type: "pot.instance.small"},
			err: regexp.MustCompile(`(?i)name`),
		},
		// invalid type
		{
			in:  model.Resource{Name: "my foo", Type: "foo.bar.baz"},
			err: regexp.MustCompile(`(?i)type`),
		},
	}

	for _, tt := range tests {
		func() {
			r, td := mustInitRouter(true)
			defer td()

			w := mustRequest(t, r, http.MethodPost, "/resources", tt.in)
			if w == nil {
				return
			}
			if !assert.Equal(t, http.StatusBadRequest, w.Code) {
				return // unmarshalling will fail
			}

			var res map[string]interface{}
			mustDecode(t, w, &res)
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
		w := mustRequest(t, r, http.MethodGet, fmt.Sprintf("/resources/%d", tt.userID))
		if w == nil {
			continue
		}
		if !assert.Equal(t, tt.code, w.Code) {
			continue
		}
		if tt.code != http.StatusOK {
			continue
		}

		var res []*model.Resource
		mustDecode(t, w, &res)
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
		w := mustRequest(t, r, http.MethodGet, fmt.Sprintf("/resources/%d/%d", tt.userID, tt.id))
		if w == nil {
			continue
		}
		if !assert.Equal(t, tt.code, w.Code) {
			continue
		}
		if w.Code != http.StatusOK {
			continue
		}

		var rc model.Resource
		mustDecode(t, w, &rc)
		assert.NotEmpty(t, rc)
		assert.Equal(t, tt.id, rc.ID)
		assert.Equal(t, tt.userID, rc.UserID)
		assert.NotEmpty(t, rc.Name)
		assert.NotEmpty(t, rc.Type)
	}
}

func TestDeleteForUser(t *testing.T) {
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
		w := mustRequest(t, r, http.MethodDelete, fmt.Sprintf("/resources/%d/%d", tt.userID, tt.id))
		if w == nil {
			continue
		}
		if !assert.Equal(t, tt.code, w.Code) {
			continue
		}
		if w.Code != http.StatusOK {
			continue
		}

		var rc model.Resource
		mustDecode(t, w, &rc)
		assert.NotEmpty(t, rc)
		assert.Equal(t, tt.id, rc.ID)
		assert.Equal(t, tt.userID, rc.UserID)
		assert.NotEmpty(t, rc.Name)
		assert.NotEmpty(t, rc.Type)

		// test if resource was really deleted
		w = mustRequest(t, r, http.MethodGet, fmt.Sprintf("/resources/%d/%d", tt.userID, tt.id))
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
		w := mustRequest(t, r, http.MethodPatch, fmt.Sprintf("/resources/%d/%d", tt.userID, tt.id), tt.in)
		if w == nil {
			continue
		}
		if !assert.Equal(t, tt.code, w.Code) {
			continue
		}
		if w.Code != http.StatusOK {
			continue
		}

		var rc model.Resource
		mustDecode(t, w, &rc)
		assert.NotEmpty(t, rc)
		assert.Equal(t, tt.id, rc.ID)
		assert.Equal(t, tt.userID, rc.UserID)
		assert.NotEmpty(t, rc.Name)
		assert.NotEmpty(t, rc.Type)
	}
}
