package main

import(
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/agent"
	"time"
	"fmt"
	"os/signal"
	"syscall"
	"context"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg := &config{}
	err := parseFlags(cfg)
	if err != nil { panic(fmt.Errorf("failed to parse flags: %w", err)) }
	if err := run(cfg); err != nil {
		log.Debug().Err(err).Msg("server error")
		panic(err)
	}
}

func run(cfg *config) (err error) {
	ctx, cancelFunc := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancelFunc()

	a := agent.NewAgent(cfg.ReportAddr)
	errCh := make(chan error, 2)

	go a.UpdateWorker.Run(ctx, errCh, time.Duration(cfg.PollInterval) * time.Second)
	go a.ReportWorker.Run(ctx, errCh, time.Duration(cfg.ReportInterval) * time.Second)

	fmt.Println("Agent working, for exit press Ctrl+C")

	for {
		select {
			case <-ctx.Done(): return nil
			case err = <-errCh: return fmt.Errorf("worker error: %w", err)
		}
	}
}
