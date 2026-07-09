package agent

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/client"
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/stats"
)

func TestReportMetrics(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(""))
	}))
	defer ts.Close()
	c := client.NewClient(ts.URL, "", []byte{})
	s := stats.NewStats()
	a := NewAgent(c, s, 1, 2, 1)
	a.reportMetrics()
}
