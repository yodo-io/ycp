package v1

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/yodo-io/ycp/pkg/model"
)

// Init router for HTTP tests. Calls v1.Setup() to make sure all routes are registered
// In case of any errors this will panic, so it's not intended for use outside of test code
// Clients must defer the returned teardown function to execute any shutdown functionality
func mustInitRouter(sampleData bool) (*gin.Engine, func()) {
	gin.SetMode(gin.TestMode)

	db := model.MustInitTestDB(sampleData)

	g := gin.New()
	g.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{}) // make sure we return json
	})

	rg := g.Group("/")
	Setup(rg, db)

	teardown := func() {
		db.Close()
	}

	return g, teardown
}

// Decode JSON payload from buffer into receiver value
func mustDecode(t *testing.T, w *httptest.ResponseRecorder, to interface{}) bool {
	d, _ := ioutil.ReadAll(w.Body)
	if err := json.Unmarshal(d, to); err != nil {
		t.Fatal(err)
		return false
	}
	return true
}

// Check if status meets expectations, if not, abort test
func mustBeStatus(t *testing.T, r string, exp int, w *httptest.ResponseRecorder) bool {
	if exp != w.Code {
		t.Fatalf("Expected %s to respond %d, got %d", r, exp, w.Code)
		return false
	}
	return true
}

// Check if status meets expectations, if not, register an error but don't abort
func shouldBeStatus(t *testing.T, r string, exp int, w *httptest.ResponseRecorder) bool {
	if exp != w.Code {
		t.Errorf("Expected %s to respond %d, got %d", r, exp, w.Code)
		return false
	}
	return true
}

// Do request using httptest.ResponseRecorder, optionally send payload. Only a single value is
// considered for payload, additional values are ignored.
func mustRequest(t *testing.T, r *gin.Engine, method string, path string, payload ...interface{}) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	var data io.Reader

	if len(payload) == 1 {
		b, err := json.Marshal(payload[0])
		if err != nil {
			t.Fatal(err)
			return nil
		}
		data = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, path, data)
	if err != nil {
		t.Fatal(err)
		return nil
	}
	req.Header.Add("Content-type", "application/json")

	r.ServeHTTP(w, req)
	return w
}
