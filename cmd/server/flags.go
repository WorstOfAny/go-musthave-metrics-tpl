package main

import(
	"flag"
	"github.com/caarlos0/env/v11"
)

var flagRunAddr string
var flagStoreInterval int
var flagFileStoragePath string
var flagRestoreStorage bool

type config struct {
	RunAddr *string `env:"ADDRESS"`
	StoreInterval *int `env:"STORE_INTERVAL"`
	FileStoragePath *string `env:"FILE_STORAGE_PATH"`
	RestoreStorage *bool `env:"RESTORE"`
}

func parseFlags() {

	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&flagStoreInterval, "i", 300, "Write storage to restore file interval in seconds")
	flag.StringVar(&flagFileStoragePath, "f", "store.json", "Storage restore file")
	flag.BoolVar(&flagRestoreStorage, "r", false, "Load from storage restore file")
	flag.Parse()

	var cfg config
	_ = env.Parse(&cfg)
	if cfg.RunAddr != nil { flagRunAddr = *cfg.RunAddr }
	if cfg.StoreInterval != nil { flagStoreInterval = *cfg.StoreInterval }
	if cfg.FileStoragePath != nil { flagFileStoragePath = *cfg.FileStoragePath }
	if cfg.RestoreStorage != nil { flagRestoreStorage = *cfg.RestoreStorage }
}
