package agent

import(
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/stats"
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/client"
	"time"
	"path"
	"context"
	"net/url"
)

type agent struct {
	stats *stats.Stats
	client *client.Client
	reportURL *url.URL
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
}

func (w *worker) Run(ctx context.Context, errCh chan error, delay time.Duration) (err error) {
	for {
		select {
			case <-ctx.Done(): return
			case <-time.After(delay):
				err = w.action()
				if err != nil { errCh <- err }
		}
	}
}

func (a *agent) reportMetrics() (err error) {
	for metric := range a.stats.AllMetrics() {
		err = a.client.Post(a.reportURL.String() + "/" + path.Join("update", metric.MType, metric.ID, metric.StringValue()))
	}
	*a.stats.PollCount = 0
	return err
}
