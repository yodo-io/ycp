package v1

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/yodo-io/ycp/pkg/api/test"
	"github.com/yodo-io/ycp/pkg/model"
)

func TestGetQuotaForUser(t *testing.T) {
	r, td := mustInitRouter(true)
	defer td()

	tests := []struct {
		userID uint
		len    int
		code   int
	}{
		{
			userID: 1,
			len:    1,
			code:   http.StatusOK,
		},
		{
			userID: 2,
			len:    0,
			code:   http.StatusOK,
		},
		{
			userID: 20,
			code:   http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		w := test.MustRecord(t, r, http.MethodGet, fmt.Sprintf("/quotas/%d", tt.userID))
		if w == nil {
			continue
		}
		if !assert.Equal(t, w.Code, tt.code) {
			continue
		}
		if w.Code != http.StatusOK {
			continue
		}

		var res []model.Quota
		test.MustDecode(t, w, &res)

		assert.Len(t, res, tt.len)
		for _, q := range res {
			assert.NotEmpty(t, q.Type)
			assert.Equal(t, tt.userID, q.UserID)
		}
	}
}

func TestCreateQuota(t *testing.T) {
	r, td := mustInitRouter(true)
	defer td()

	tests := []struct {
		userID uint
		in     model.Quota
		out    model.Quota
		code   int
	}{
		// ok
		{
			userID: 1,
			in:     model.Quota{Type: "pot.instance.small", Value: 20},
			out:    model.Quota{Type: "pot.instance.small", Value: 20, UserID: 1},
			code:   http.StatusCreated,
		},
		// should fix userid mismatch
		{
			userID: 1,
			in:     model.Quota{Type: "pot.instance.small", Value: 20, UserID: 10},
			out:    model.Quota{Type: "pot.instance.small", Value: 20, UserID: 1},
			code:   http.StatusCreated,
		},
		// non-existing resource
		{
			userID: 1,
			in:     model.Quota{Type: "pitchfork.instance.3s", Value: 5},
			code:   http.StatusBadRequest,
		},
		// non-existing user
		{
			userID: 10,
			in:     model.Quota{Type: "pot.instance.small"},
			code:   http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		// t.Logf("%#v\n", tt.in)
		w := test.MustRecord(t, r, http.MethodPost, fmt.Sprintf("/quotas/%d", tt.userID), tt.in)
		if w == nil {
			continue
		}
		if !assert.Equal(t, tt.code, w.Code) {
			continue
		}
		if w.Code != http.StatusCreated {
			continue
		}

		var res model.Quota
		test.MustDecode(t, w, &res)
		assert.NotZero(t, res.ID)
		assert.Equal(t, tt.in.Type, res.Type)
		assert.Equal(t, tt.in.Value, res.Value)
		assert.Equal(t, tt.userID, res.UserID)
	}
}

func TestDeleteQuota(t *testing.T) {
	r, td := mustInitRouter(true)
	defer td()

	tests := []struct {
		id     uint
		userID uint
		code   int
	}{
		{userID: 1, id: 1, code: http.StatusOK},
		{userID: 2, id: 1, code: http.StatusNotFound},
		{userID: 1, id: 3, code: http.StatusNotFound},
	}

	for _, tt := range tests {
		w := test.MustRecord(t, r, http.MethodDelete, fmt.Sprintf("/quotas/%d/%d", tt.userID, tt.id))
		if w == nil {
			continue
		}
		if !assert.Equal(t, tt.code, w.Code) {
			continue
		}
		if w.Code != http.StatusOK {
			continue
		}

		var q model.Quota
		test.MustDecode(t, w, &q)
		assert.NotEmpty(t, q)
		assert.Equal(t, tt.id, q.ID)
		assert.Equal(t, tt.userID, q.UserID)
		assert.NotEmpty(t, q.Type)

		// test if resource was really deleted
		w = test.MustRecord(t, r, http.MethodGet, fmt.Sprintf("/quotas/%d/%d", tt.userID, tt.id))
		assert.Equal(t, http.StatusNotFound, w.Code)
	}
}

func TestUpdateQuota(t *testing.T) {
	r, td := mustInitRouter(true)
	defer td()

	tests := []struct {
		id     uint
		userID uint
		value  int
		code   int
	}{
		{
			id:     1,
			userID: 1,
			value:  100,
			code:   http.StatusOK,
		},
		{
			id:     1,
			userID: 1,
			value:  0,
			code:   http.StatusOK,
		},
		{
			id:     20,
			userID: 1,
			value:  100,
			code:   http.StatusNotFound,
		},
		{
			id:     1,
			userID: 20,
			value:  100,
			code:   http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		in := gin.H{"value": tt.value}
		w := test.MustRecord(t, r, http.MethodPatch, fmt.Sprintf("/quotas/%d/%d", tt.userID, tt.id), in)

		if w == nil {
			continue
		}
		if !assert.Equal(t, tt.code, w.Code) {
			continue
		}
		if w.Code != http.StatusOK {
			continue
		}

		var q model.Quota
		test.MustDecode(t, w, &q)
		assert.NotEmpty(t, q)
		assert.Equal(t, tt.id, q.ID)
		assert.Equal(t, tt.userID, q.UserID)
		assert.Equal(t, tt.value, q.Value)
		assert.NotEmpty(t, q.Type)
	}
}
