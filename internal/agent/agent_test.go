package agent

import(
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
	"context"
	"os/signal"
	"syscall"
	"errors"
)

func TestRun(t *testing.T) {
	type want struct {
		result int
	}

	testcases := []struct{
		name string
		sign syscall.Signal
		err error
	}{
		{
			name: "SIGINT",
			sign: syscall.SIGINT,
		},
		{
			name: "SIGTERM",
			sign: syscall.SIGTERM,
		},
		{
			name: "With error",
			sign: syscall.SIGTERM,
			err: errors.New("test error"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			callTimes := 0
			errCh := make(chan error, 1)

			ctx, cancelFunc := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer cancelFunc()

			worker := worker{action: func() error { callTimes++; return tc.err }}
			go worker.Run(ctx, errCh, 1 * time.Second)
			time.Sleep(4 * time.Second + 5 * time.Millisecond)

			err := syscall.Kill(syscall.Getpid(), tc.sign)
			if err != nil {
				t.Fatal(err)
			}

			select {
				case err := <- errCh:
					if tc.err != nil {
						assert.Equal(t, tc.err, err)
					} else {
						t.Error("ошибка не ожидалась")
					}
				case <- ctx.Done():
					assert.Equal(t, context.Canceled, ctx.Err())
					assert.Equal(t, 4, callTimes)
				case <- time.After(2 * time.Second):
					t.Error("контекст не реагирует на сигнал")
			}
		})
	}
}

func TestReportMetrics(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(""))
	}))
	defer ts.Close()
	reportURL, _ := url.Parse(ts.URL)
	a := NewAgent(reportURL.Host)
	a.reportMetrics()
}
