package main

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v11"
)

type config struct {
	ReportAddr     string `env:"ADDRESS"`
	Key            string `env:"KEY"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	RateLimit      int    `env:"RATE_LIMIT"`
}

func parseFlags(cfg *config) (err error) {
	flag.StringVar(&cfg.ReportAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&cfg.Key, "k", "", "key")
	flag.IntVar(&cfg.PollInterval, "p", 2, "poll interval")
	flag.IntVar(&cfg.ReportInterval, "r", 10, "report interval")
	flag.IntVar(&cfg.RateLimit, "l", 10, "requests rate limit")
	flag.Parse()

	err = env.Parse(cfg)
	if err != nil {
		return fmt.Errorf("failed to read configuration from environment variables: %w", err)
	}

	return nil
}
