package main

import(
	"flag"
	"github.com/caarlos0/env/v11"
	"fmt"
)

//type dbConfig struct {
//	Host *string `env:"HOST"`
//	Port *int `env:"PORT"`
//	User *string `env:"USER"`
//	Password *string `env:"PASSWORD"`
//	Name *string `env:"NAME"`
//}

type config struct {
	RunAddr string `env:"ADDRESS"`
	StoreInterval int `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	RestoreStorage bool `env:"RESTORE"`
	DatabaseDSN string `env:"DATABASE_DSN"`
	//DB dbConfig `envPrefix:"DATABASE_"`
}

var cfg config

func parseFlags() (err error) {
	cfg = config{}
	flag.StringVar(&cfg.RunAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&cfg.StoreInterval, "i", 300, "Write storage to restore file interval in seconds")
	flag.StringVar(&cfg.FileStoragePath, "f", "store.json", "Storage restore file")
	flag.BoolVar(&cfg.RestoreStorage, "r", false, "Load from storage restore file")
	flag.StringVar(&cfg.DatabaseDSN, "d", "", "Database data source name")
	flag.Parse()

	fmt.Println(cfg)

	err = env.Parse(&cfg)
	if err != nil {
		return fmt.Errorf("failed to read configuration from environment variables: %w", err)
	}
	fmt.Println(cfg)

	return nil
}
