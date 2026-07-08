package main

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v11"

	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/repository"
)

type config struct {
	RunAddr    string `env:"ADDRESS"`
	Key        string `env:"KEY"`
	AuditURL   string `env:"AUDIT_FILE"`
	AuditFile  string `env:"AUDIT_URL"`
	RepoConfig repository.Config
}

func parseFlags(cfg *config) (err error) {
	flag.StringVar(&cfg.RunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&cfg.Key, "k", "", "key")
	flag.StringVar(&cfg.AuditURL, "audit-url", "", "audit log url")
	flag.StringVar(&cfg.AuditFile, "audit-file", "", "audit log file")
	flag.IntVar(&cfg.RepoConfig.StoreInterval, "i", 300, "Write storage to restore file interval in seconds")
	flag.StringVar(&cfg.RepoConfig.FileStoragePath, "f", "store.json", "Storage restore file")
	flag.BoolVar(&cfg.RepoConfig.RestoreStorage, "r", false, "Load from storage restore file")
	flag.StringVar(&cfg.RepoConfig.DatabaseDSN, "d", "", "Database data source name")
	flag.Parse()

	err = env.Parse(cfg)
	if err != nil {
		return fmt.Errorf("failed to read configuration from environment variables: %w", err)
	}

	return nil
}
