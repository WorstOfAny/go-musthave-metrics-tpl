package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v11"

	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/repository"
)

type config struct {
	RunAddr    string            `env:"ADDRESS" json:"address"`
	Key        string            `env:"KEY" json:"key"`
	AuditURL   string            `env:"AUDIT_FILE" json:"audit_url"`
	AuditFile  string            `env:"AUDIT_URL" json:"audit_file"`
	RepoConfig repository.Config `json:"repository"`
	SecretPath string            `env:"SECRET_PATH" json:"crypto_key"`
}

func newConfig() *config {
	return &config{
		RunAddr: "localhost:8080",
		RepoConfig: repository.Config{
			StoreInterval:   300 * time.Second,
			FileStoragePath: "store.json",
		},
	}
}

func parseFlags(cfg *config) (err error) {
	var (
		confJSONPath  *string
		serverAddr    *string
		key           *string
		auditURL      *string
		auditFile     *string
		cryptoKey     *string
		storeInterval *time.Duration
		storeFile     *string
		storeRestore  *bool
		databaseDSN   *string
	)
	confJSONPath = flag.String("config", "", "json config file path")
	flag.StringVar(confJSONPath, "c", "", "json config file path")
	serverAddr = flag.String("a", "localhost:8080", "address and port to run server")
	key = flag.String("k", "", "key")
	auditURL = flag.String("audit-url", "", "audit log url")
	auditFile = flag.String("audit-file", "", "audit log file")
	cryptoKey = flag.String("crypto-key", "", "private key path")
	storeInterval = flag.Duration("i", 300*time.Second, "Write storage to restore file interval in seconds")
	storeFile = flag.String("f", "store.json", "Storage restore file")
	storeRestore = flag.Bool("r", false, "Load from storage restore file")
	databaseDSN = flag.String("d", "", "Database data source name")
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
			cfg.RunAddr = *serverAddr
		case "k":
			cfg.Key = *key
		case "audit-url":
			cfg.AuditURL = *auditURL
		case "audit-file":
			cfg.AuditFile = *auditFile
		case "crypto-key":
			cfg.SecretPath = *cryptoKey
		case "i":
			cfg.RepoConfig.StoreInterval = *storeInterval
		case "f":
			cfg.RepoConfig.FileStoragePath = *storeFile
		case "r":
			cfg.RepoConfig.RestoreStorage = *storeRestore
		case "d":
			cfg.RepoConfig.DatabaseDSN = *databaseDSN
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
