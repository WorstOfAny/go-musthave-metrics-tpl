package main

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"

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

	cfg := &config{}
	err := parseFlags(cfg)
	if err != nil {
		panic(fmt.Errorf("failed to parse flags: %w", err))
	}
	if err := run(cfg); err != nil {
		log.Debug().Err(err).Msg("server error")
		panic(err)
	}
}

func run(cfg *config) (err error) {
	ctx, cancelFunc := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
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

	return nil
}
