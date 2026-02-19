package agent

import(
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/stats"
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/client"
	"time"
	"path"
	"context"
	"net/url"
)

var statsObject = stats.NewStats()
var clientObject = client.NewClient()

var ReportURL = &url.URL{ Scheme: "http", Host: "localhost:8080" }

type Worker struct {
	action func()
}


var ReportWorker = &Worker{ action: reportMetrics }
var UpdateWorker = &Worker{ action: statsObject.Update }

func (w *Worker) Run(ctx context.Context, delay time.Duration) {
	for {
		select {
			case <-ctx.Done(): return
			case <-time.After(delay): w.action()
		}
	}
}

func reportMetrics() {
	for metric := range statsObject.AllMetrics() {
		clientObject.Post(ReportURL.String() + "/" + path.Join("update", metric.MType, metric.ID, metric.StringValue()))
	}
	*statsObject.PollCount = 0
}
