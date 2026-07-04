package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/client"
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/stats"
)

type agent struct {
	stats       *stats.Stats
	client      *client.Client
	rateLimit   int
	reportDelay time.Duration
	pollDelay   time.Duration
	mu          sync.Mutex
}

func NewAgent(client *client.Client, stats *stats.Stats, ratelimit int, reportInterval int, pollInterval int) *agent {
	return &agent{
		stats:       stats,
		client:      client,
		rateLimit:   ratelimit,
		reportDelay: time.Duration(reportInterval) * time.Second,
		pollDelay:   time.Duration(pollInterval) * time.Second,
	}
}

func (a *agent) Start(ctx context.Context) error {
	g, errGrCtx := errgroup.WithContext(ctx)
	reportJobs := jobsGenerator(errGrCtx, a.reportMetrics, a.reportDelay)
	updateRTJobs := jobsGenerator(errGrCtx, a.stats.UpdateRT, a.pollDelay)
	updateVMJobs := jobsGenerator(errGrCtx, a.stats.UpdateVM, a.pollDelay)

	for i := 1; i <= a.rateLimit; i++ {
		worker(g, reportJobs)
	}
	worker(g, updateRTJobs)
	worker(g, updateVMJobs)

	fmt.Println("Agent working, for exit press Ctrl+C")

	if err := g.Wait(); err != nil {
		return fmt.Errorf("worker error: %w", err)
	}

	return nil
}

func jobsGenerator(ctx context.Context, job func() error, delay time.Duration) chan func() error {
	jobs := make(chan func() error)
	go func() {
		defer close(jobs)
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(delay):
				jobs <- job
			}
		}
	}()

	return jobs
}

func worker(g *errgroup.Group, jobs <-chan func() error) {
	g.Go(func() error {
		for j := range jobs {
			err := j()
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (a *agent) reportMetrics() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	body, err := json.Marshal(a.stats)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	err = a.client.Post("/updates", body)
	if err != nil {
		return fmt.Errorf("failed to send data to server: %w", err)
	}

	*a.stats.PollCount = 0
	return nil
}
