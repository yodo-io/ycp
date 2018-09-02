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

func TestGetUsers(t *testing.T) {
	r, td := mustInitRouter(true)
	defer td()

	w, _ := doRequest(r, http.MethodGet, "/users", nil)
	var users []model.User

	mustBeStatus(t, "GET /users", http.StatusOK, w)
	mustDecode(t, w, &users)
	assert.NotEmpty(t, users, "Expected some users")
}

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

			w, _ := doRequest(r, http.MethodPost, "/users", tt.in)

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

func TestValidaton(t *testing.T) {

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

			w, _ := doRequest(r, http.MethodPost, "/users", tt.in)
			var res map[string]interface{}

			if !assert.Equal(t, http.StatusBadRequest, w.Code) {
				return // unmarshalling will fail
			}
			mustDecode(t, w, &res)
			assert.NotEmpty(t, res["error"])
			assert.Regexp(t, tt.err, res["error"].(string))
		}()
	}
}
