package client

import(
	"net/http/httptest"
	"net/http"
	"testing"
	"compress/gzip"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	NewClient("", "")
}

func TestPost(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gz.Write([]byte("{\"a\": 1}"))
	}))
	c := NewClient(ts.URL, "")

	err := c.Post("", []byte("{}"))
	assert.NoError(t, err)

	ts.Close()

	err = c.Post(ts.URL, []byte("{}"))
	assert.Error(t, err)
}
