package main

import(
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/agent"
	"time"
	"context"
)

const(
	pollInterval = 2 * time.Second
	reportInterval = 5 * time.Second
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	go agent.UpdateWorker.Run(ctx, pollInterval)
	go agent.ReportWorker.Run(ctx, reportInterval)
	<-ctx.Done()
	return ctx.Err()
}
