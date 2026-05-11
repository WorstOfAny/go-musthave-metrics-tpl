package agent

import(
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/stats"
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/client"
	"time"
	"context"
	"net/url"
	"encoding/json"
	"sync"
	"github.com/rs/zerolog/log"
	"fmt"
)

type agent struct {
	stats *stats.Stats
	client *client.Client
	reportURL *url.URL
	mu sync.Mutex
	UpdateWorker *worker
	ReportWorker *worker
}

func NewAgent(reportAddr string) *agent {
	a := agent{
		stats: stats.NewStats(),
		client: client.NewClient(),
		reportURL: &url.URL{ Scheme: "http", Host: reportAddr },
	}

	a.UpdateWorker = &worker{ action: a.stats.Update }
	a.ReportWorker = &worker{ action: a.reportMetrics }

	return &a
}

type worker struct {
	action func() error
	mu sync.Mutex
}

func (w *worker) Run(ctx context.Context, errCh chan error, delay time.Duration) (err error) {
	for {
		select {
			case <-ctx.Done(): return
			case <-time.After(delay):
				w.mu.Lock()
				err = w.action()
				if err != nil {
					log.Debug().Err(err).Msg("worker action err")
					errCh <- fmt.Errorf("failed to execute worker action: %w", err)
					return
				}
				w.mu.Unlock()
		}
	}
}

func (a *agent) reportMetrics() (error) {
	a.mu.Lock()

	body, err := json.Marshal(a.stats)
	if err != nil { return fmt.Errorf("failed to marshal data: %w", err) }

	err = a.client.Post(a.reportURL.String() + "/updates", body)
	if err != nil {
		return fmt.Errorf("failed to send data to server: %w", err)
	}

	*a.stats.PollCount = 0
	a.mu.Unlock()
	return nil
}
