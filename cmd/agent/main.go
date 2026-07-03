package main

import(
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/agent"
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/stats"
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/client"
	"fmt"
	"os/signal"
	"syscall"
	"context"
	"github.com/rs/zerolog/log"
	"net/url"
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

	s := stats.NewStats()
	c := client.NewClient((&url.URL{ Scheme: "http", Host: cfg.ReportAddr }).String(), cfg.Key)

	a := agent.NewAgent(c, s, cfg.RateLimit, cfg.ReportInterval, cfg.PollInterval)
	err = a.Start(ctx)

	if err != nil {
		return fmt.Errorf("agent error: %w", err)
	}

	return nil
}
