package agent

import(
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
	"context"
)

func TestRun(t *testing.T) {
	callTimes := 0

	ctx, cancelFunc := context.WithCancel(context.Background())

	worker := Worker{action: func() { callTimes++ }}

	go worker.Run(ctx, 1 * time.Second)
	time.Sleep(4 * time.Second + 5 * time.Millisecond)
	cancelFunc()

	assert.Equal(t, 4, callTimes)
}

func TestReportMetrics(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "text/plain", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(""))
	}))
	defer ts.Close()

	ReportURL, _ = url.Parse(ts.URL)
	reportMetrics()
}
