package main

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/agent"
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/client"
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/stats"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Build version: %s\n", buildVersion)
	fmt.Fprintf(&buf, "Build date: %s\n", buildDate)
	fmt.Fprintf(&buf, "Build commit: %s\n", buildCommit)
	os.Stdout.Write(buf.Bytes())

	cfg := newConfig()
	err := parseFlags(cfg)
	if err != nil {
		panic(fmt.Errorf("failed to parse flags: %w", err))
	}
	if err := run(cfg); err != nil {
		if errors.Is(err, context.Canceled) {
			log.Info().Msg("agent gracefully shutdown")
			return
		}
		log.Error().Err(err).Msg("agent error")
		panic(err)
	}
}

func run(cfg *config) (err error) {
	ctx, cancelFunc := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancelFunc()

	s := stats.NewStats()
	certBytes, err := os.ReadFile(cfg.SecretPath)
	if err != nil {
		return fmt.Errorf("failed read server public key")
	}
	c := client.NewClient((&url.URL{Scheme: "http", Host: cfg.ReportAddr}).String(), cfg.Key, certBytes)

	a := agent.NewAgent(c, s, cfg.RateLimit, cfg.ReportInterval, cfg.PollInterval)
	err = a.Start(ctx)

	if err != nil {
		return fmt.Errorf("agent error: %w", err)
	}

	<-ctx.Done()

	if ctx.Err() != nil {
		return ctx.Err()
	}

	return nil
}
