package main

import(
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/agent"
	"time"
	"context"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	parseFlags()
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	agent.ReportURL.Host = flagReportAddr
	go agent.UpdateWorker.Run(ctx, time.Duration(flagPollInterval) * time.Second)
	go agent.ReportWorker.Run(ctx, time.Duration(flagReportInterval) * time.Second)
	<-ctx.Done()
	return ctx.Err()
}
