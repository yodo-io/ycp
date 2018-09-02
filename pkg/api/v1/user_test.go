package v1

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yodo-io/ycp/pkg/model"
)

func TestCreateUser(t *testing.T) {
	tests := []struct {
		in  model.User
		out model.User
	}{
		{
			in:  model.User{Email: "john@example.org", Password: "pass"},
			out: model.User{ID: 1, Email: "john@example.org", Role: "user"},
		},
		{
			in:  model.User{Email: "john@example.org", Password: "pass", Role: "admin"},
			out: model.User{ID: 1, Email: "john@example.org", Role: "admin"},
		},
	}

	for _, tt := range tests {
		func() {
			r, td := mustInitRouter(false)
			defer td()

			w := mustRequest(t, r, http.MethodPost, "/users", tt.in)
			if w == nil {
				return
			}

			if !assert.Equal(t, http.StatusCreated, w.Code) {
				d, _ := ioutil.ReadAll(w.Body)
				fmt.Println(string(d))
				return // unmarshalling will fail
			}

			var res model.User
			mustDecode(t, w, &res)
			assert.Equal(t, tt.out, res)
		}()
	}
}

func TestValidation(t *testing.T) {

	tests := []struct {
		in  model.User
		err *regexp.Regexp
	}{
		// password required
		{
			in:  model.User{Email: "john@example.org"},
			err: regexp.MustCompile(`(?i)password`),
		},
		// email required
		{
			in:  model.User{Password: "password"},
			err: regexp.MustCompile(`(?i)email`),
		},
		// role must exist
		{
			in:  model.User{Email: "john@example.org", Password: "password", Role: "superuser"},
			err: regexp.MustCompile(`(?i)role`),
		},
		// email must be valid
		{
			in:  model.User{Email: "email", Password: "password"},
			err: regexp.MustCompile(`(?i)email`),
		},
	}

	for _, tt := range tests {
		func() {
			r, td := mustInitRouter(false)
			defer td()

			w := mustRequest(t, r, http.MethodPost, "/users", tt.in)
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

func TestGetUsers(t *testing.T) {
	r, td := mustInitRouter(true)
	defer td()

	w := mustRequest(t, r, http.MethodGet, "/users")
	if w == nil {
		return
	}
	if !assert.Equal(t, http.StatusOK, w.Code) {
		return
	}

	var res []model.User
	mustDecode(t, w, &res)
	assert.NotEmpty(t, res)

	for _, u := range res {
		assert.NotEmpty(t, u.Email)
		assert.NotEmpty(t, u.Role)
		assert.Empty(t, u.Password)
	}
}

func TestGetUser(t *testing.T) {
	r, td := mustInitRouter(true)
	defer td()

	tests := []struct {
		id   uint
		code int
	}{
		{id: 1, code: http.StatusOK},
		{id: 20, code: http.StatusNotFound},
	}

	for _, tt := range tests {
		w := mustRequest(t, r, http.MethodGet, fmt.Sprintf("/users/%d", tt.id))
		if w == nil {
			continue
		}
		if !assert.Equal(t, tt.code, w.Code) {
			continue
		}

		if w.Code != http.StatusOK {
			continue
		}

		var u model.User
		mustDecode(t, w, &u)
		assert.NotEmpty(t, u)
		assert.Equal(t, tt.id, u.ID)
		assert.NotEmpty(t, u.Email)
		assert.NotEmpty(t, u.Role)
		assert.Empty(t, u.Password)
	}
}

func TestDeleteUser(t *testing.T) {
	r, td := mustInitRouter(true)
	defer td()

	tests := []struct {
		id   uint
		code int
	}{
		{id: 1, code: http.StatusOK},
		{id: 20, code: http.StatusNotFound},
	}

	for _, tt := range tests {
		w := mustRequest(t, r, http.MethodDelete, fmt.Sprintf("/users/%d", tt.id))
		if w == nil {
			continue
		}
		if !assert.Equal(t, tt.code, w.Code) {
			continue
		}
		if w.Code != http.StatusOK {
			continue
		}

		var u model.User
		mustDecode(t, w, &u)
		assert.NotEmpty(t, u)
		assert.Equal(t, tt.id, u.ID)
		assert.NotEmpty(t, u.Email)
		assert.NotEmpty(t, u.Role)
		assert.Empty(t, u.Password)

		// make sure user was really deleted
		w = mustRequest(t, r, http.MethodGet, fmt.Sprintf("/users/%d", tt.id))
		assert.Equal(t, http.StatusNotFound, w.Code)
	}
}

func TestUpdateUser(t *testing.T) {
	r, td := mustInitRouter(true)
	defer td()

	tests := []struct {
		id   uint
		user model.User
		code int
	}{
		{id: 1, user: model.User{Email: "jane@acme.org"}, code: http.StatusOK},
		{id: 1, user: model.User{Role: "admin"}, code: http.StatusOK},
		{id: 20, user: model.User{Email: "jane@acme.org"}, code: http.StatusNotFound},
		{id: 1, user: model.User{Email: "jane"}, code: http.StatusBadRequest},
		{id: 1, user: model.User{Role: "foo"}, code: http.StatusBadRequest},
	}

	for _, tt := range tests {
		w := mustRequest(t, r, http.MethodPatch, fmt.Sprintf("/users/%d", tt.id), tt.user)
		if w == nil {
			continue
		}
		if !assert.Equal(t, tt.code, w.Code) {
			continue
		}
		if w.Code != http.StatusOK {
			continue
		}

		var u model.User
		mustDecode(t, w, &u)
		assert.NotEmpty(t, u)
		assert.Equal(t, tt.id, u.ID)
		assert.NotEmpty(t, u.Email)
		assert.NotEmpty(t, u.Role)
	}
}
