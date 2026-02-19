package client

import(
	"net/http/httptest"
	"net/http"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	NewClient()
}

func TestPost(t *testing.T) {
	c := NewClient()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "text/plain", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(""))
	}))
	defer ts.Close()	

	c.Post(ts.URL)
}
