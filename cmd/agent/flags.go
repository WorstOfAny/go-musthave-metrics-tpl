package main

import(
	"flag"
	"github.com/caarlos0/env/v11"
	"fmt"
)

type config struct {
	ReportAddr string `env:"ADDRESS"`
	PollInterval int `env:"POLL_INTERVAL"`
	ReportInterval int `env:"REPORT_INTERVAL"`
}

var cfg config

func parseFlags() (err error) {
	cfg = config{}
	flag.StringVar(&cfg.ReportAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&cfg.PollInterval, "p", 2, "poll interval")
	flag.IntVar(&cfg.ReportInterval, "r", 10, "report interval")
	flag.Parse()

	err = env.Parse(&cfg)
	if err != nil {
		return fmt.Errorf("failed to read configuration from environment variables: %w", err)
	}

	return nil
}
