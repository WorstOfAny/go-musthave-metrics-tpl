package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v11"
)

type config struct {
	ReportAddr     string        `env:"ADDRESS" json:"report_address"`
	Key            string        `env:"KEY" json:"key"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" json:"poll_interval"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" json:"report_interval"`
	RateLimit      int           `env:"RATE_LIMIT" json:"rate_limit"`
	SecretPath     string        `env:"SECRET_PATH" json:"crypto_key"`
}

func (c *config) UnmarshalJSON(data []byte) error {
	type Alias config
	al := &struct {
		PollInterval   string `json:"poll_interval"`
		ReportInterval string `json:"report_interval"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, al); err != nil {
		return err
	}

	pollInterval, err := time.ParseDuration(al.PollInterval)
	if err != nil {
		return fmt.Errorf("failed parse duration: %w", err)
	}

	reportInterval, err := time.ParseDuration(al.ReportInterval)
	if err != nil {
		return fmt.Errorf("failed parse duration: %w", err)
	}

	c.PollInterval = pollInterval
	c.ReportInterval = reportInterval
	return nil
}

func newConfig() *config {
	return &config{
		ReportAddr:     "localhost:8080",
		PollInterval:   2 * time.Second,
		ReportInterval: 10 * time.Second,
		RateLimit:      10,
	}
}

func parseFlags(cfg *config) (err error) {
	var (
		confJSONPath   *string
		reportAddr     *string
		key            *string
		pollInterval   *time.Duration
		reportInterval *time.Duration
		cryptoKey      *string
		rateLimit      *int
	)
	confJSONPath = flag.String("config", "", "json config file path")
	flag.StringVar(confJSONPath, "c", "", "json config file path")
	reportAddr = flag.String("a", "localhost:8080", "address and port to run server")
	key = flag.String("k", "", "key")
	pollInterval = flag.Duration("p", 2*time.Second, "poll interval")
	reportInterval = flag.Duration("r", 10*time.Second, "report interval")
	rateLimit = flag.Int("l", 10, "requests rate limit")
	cryptoKey = flag.String("crypto-key", "", "server public key path")
	flag.Parse()

	flag.Visit(func(f *flag.Flag) {
		if f.Name == "c" || f.Name == "config" {
			readAndParseConfig(*confJSONPath, cfg)
		}
	})

	val, ok := os.LookupEnv("CONFIG")

	if ok {
		readAndParseConfig(val, cfg)
	}

	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "a":
			cfg.ReportAddr = *reportAddr
		case "k":
			cfg.Key = *key
		case "p":
			cfg.PollInterval = *pollInterval
		case "r":
			cfg.ReportInterval = *reportInterval
		case "l":
			cfg.RateLimit = *rateLimit
		case "crypto-key":
			cfg.SecretPath = *cryptoKey
		}
	})

	err = env.Parse(cfg)
	if err != nil {
		return fmt.Errorf("failed to read configuration from environment variables: %w", err)
	}

	return nil
}

func readAndParseConfig(path string, cfg *config) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("failed to read configuration from config file: %w", err))
	}

	err = json.Unmarshal(bytes, cfg)

	if err != nil {
		panic(fmt.Errorf("failed to unmarshal data from configuration file: %w", err))
	}
}
