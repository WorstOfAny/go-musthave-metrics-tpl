package main

import(
	"flag"
	"github.com/caarlos0/env/v11"
)

var flagReportAddr string
var flagPollInterval int
var flagReportInterval int

type config struct {
	ReportAddr *string `env:"ADDRESS"`
	PollInterval *int `env:"POLL_INTERVAL"`
	ReportInterval *int `env:"REPORT_INTERVAL"`
}

func parseFlags() {
	flag.StringVar(&flagReportAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&flagPollInterval, "p", 2, "poll interval")
	flag.IntVar(&flagReportInterval, "r", 10, "report interval")
	flag.Parse()

	var cfg config
	_ = env.Parse(&cfg)
	if cfg.ReportAddr != nil { flagReportAddr = *cfg.ReportAddr }
	if cfg.PollInterval != nil { flagPollInterval = *cfg.PollInterval }
	if cfg.ReportInterval != nil { flagReportInterval = *cfg.ReportInterval }
}
