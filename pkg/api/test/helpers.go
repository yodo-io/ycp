package test

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/yodo-io/ycp/pkg/api"
)

// NewRouter creates a new Gin router in test mode
func NewRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	g := gin.New()
	g.NoRoute(api.NotFound)
	return g
}

// MustRecord will exectue request using httptest.ResponseRecorder, optionally send payload.
// Only a single value is considered for payload, additional values are ignored.
func MustRecord(t *testing.T, r *gin.Engine, method string, path string, payload ...interface{}) *httptest.ResponseRecorder {
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

// MustDecode decodes JSON payload from buffer into receiver value. If it fails, reports
// fatal to the test runner and returns false
func MustDecode(t *testing.T, w *httptest.ResponseRecorder, to interface{}) bool {
	d, _ := ioutil.ReadAll(w.Body)
	if err := json.Unmarshal(d, to); err != nil {
		t.Fatal(err)
		return false
	}
	return true
}
